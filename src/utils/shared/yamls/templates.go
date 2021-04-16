package yamls

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	jsonpatch "github.com/evanphx/json-patch"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/kubectl/pkg/scheme"
	k8syaml "sigs.k8s.io/yaml"
)

var nonNamespacedTypes = []string{"podsecuritypolicies", "namespaces", "clusterrolebindings", "clusterroles"}
var templateTypes = []string{"daemonsets", "deployments", "statefulsets"}

// YAMLTmplArguments is a wrapper around YAMLTmplValues.
type YAMLTmplArguments struct {
	Values *map[string]interface{}
}

// ExecuteTemplatedYAMLs takes a template YAML and applies the given template values to it.
func ExecuteTemplatedYAMLs(yamls []*YAMLFile, tmplValues *YAMLTmplArguments) ([]*YAMLFile, error) {
	// Execute the template on each of the YAMLs.
	executedYAMLs := make([]*YAMLFile, len(yamls))
	for i, y := range yamls {
		yamlFile := &YAMLFile{
			Name: y.Name,
		}

		if tmplValues == nil {
			yamlFile.YAML = y.YAML
		} else {
			executedYAML, err := executeTemplate(tmplValues, y.YAML)
			if err != nil {
				return nil, err
			}
			yamlFile.YAML = executedYAML
		}
		executedYAMLs[i] = yamlFile
	}

	return executedYAMLs, nil
}

func required(str string, value string) (string, error) {
	if value != "" {
		return value, nil
	}
	return "", errors.New("Value is required")
}

func executeTemplate(tmplValues *YAMLTmplArguments, tmplStr string) (string, error) {
	funcMap := sprig.TxtFuncMap()
	funcMap["required"] = required

	tmpl, err := template.New("yaml").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, tmplValues)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

// TemplateMatchFn is a function used to determine whether or not the given resource should have the template applied.
type TemplateMatchFn func(obj map[string]interface{}, resourceType string) bool

// GenerateResourceNameMatcherFn creates a matcher function for matching the resource if the resource's name matches matchValue.
func GenerateResourceNameMatcherFn(expectedName string) TemplateMatchFn {
	fn := func(obj map[string]interface{}, resourceType string) bool {
		if md, ok := obj["metadata"]; ok {
			if name, nameOk := md.(map[string]interface{})["name"]; nameOk {
				if name == expectedName {
					return true
				}
			}
		}
		return false
	}
	return fn
}

// NamespaceScopeMatcher matches the resource if the resource is contained within a namespace.
func NamespaceScopeMatcher(obj map[string]interface{}, resourceType string) bool {
	for _, r := range nonNamespacedTypes {
		if resourceType == r {
			return false
		}
	}
	return true
}

// TemplateScopeMatcher matches the resource definition contains a template for deploying other resources.
func TemplateScopeMatcher(obj map[string]interface{}, resourceType string) bool {
	for _, r := range templateTypes {
		if resourceType == r {
			return true
		}
	}
	return false
}

// GenerateContainerNameMatcherFn creates a matcher function for matching the resource if the resource has a pod
// template with a container of the given name.
func GenerateContainerNameMatcherFn(expectedName string) TemplateMatchFn {
	fn := func(obj map[string]interface{}, resourceType string) bool {
		// Check that spec.template.spec.containers[].name in the YAML matches the given name.
		spec := obj["spec"]
		if spec == nil {
			return false
		}

		tmpl := spec.(map[string]interface{})["template"]
		if tmpl == nil {
			return false
		}

		tmplSpec := tmpl.(map[string]interface{})["spec"]
		if tmplSpec == nil {
			return false
		}
		containers := tmplSpec.(map[string]interface{})["containers"]
		if containers == nil {
			return false
		}
		containersList, ok := containers.([]interface{})
		if !ok {
			return false
		}

		for _, c := range containersList {
			container, ok := c.(map[string]interface{})
			if !ok {
				continue
			}
			if container["name"] == expectedName {
				return true
			}
		}

		return false
	}

	return fn
}

// K8sTemplateOptions specifies how the templated YAML should be constructed, by specifying selectors for which resources should
// contain the template, how the placeholder should be patched in, and what that placeholder should be replaced with.
type K8sTemplateOptions struct {
	// TemplateMatcher is a function that returns whether or not the template should be applied to the resource.
	TemplateMatcher TemplateMatchFn
	// The JSON that should be patched in, with the placeholder values.
	Patch string
	// Placeholder is the string in the YAML which should be replaced with the template value.
	Placeholder string
	// TemplateValue is the template string that should replace the placeholder.
	TemplateValue string
}

// Adds the placeholder to the object's JSON, by doing a strategic merge. The strategic merge will handle whether or not
// the patch should replace or add, depending on the K8s object schema.
func addPlaceholder(opt *K8sTemplateOptions, gvk schema.GroupVersionKind, originalJSON []byte) ([]byte, error) {
	creatorObj, err := scheme.Scheme.New(gvk)
	if err != nil {
		// Strategic merge patches are not supported for non-native K8s resources (custom CRDs).
		// We will need to perform a regular JSON patch instead.
		return jsonpatch.MergePatch(originalJSON, []byte(opt.Patch))
	}

	return strategicpatch.StrategicMergePatch(originalJSON, []byte(opt.Patch), creatorObj)
}

func addPlaceholders(rm meta.RESTMapper, decodedYAML *yaml.YAMLOrJSONDecoder, tmplOpts []*K8sTemplateOptions) (string, error) {
	// Read resource into JSON.
	ext := runtime.RawExtension{}
	err := decodedYAML.Decode(&ext)
	if err != nil {
		return "", err
	}

	_, gvk, err := unstructured.UnstructuredJSONScheme.Decode(ext.Raw, nil, nil)
	if err != nil {
		return "", err
	}

	resourceType := ""
	mapping, err := rm.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err == nil {
		k8sRes := mapping.Resource
		resourceType = k8sRes.Resource
	}

	// Decode object into readable struct.
	var unstructuredOrig unstructured.Unstructured
	unstructuredOrig.Object = make(map[string]interface{})
	var unstructBlob interface{}

	err = json.Unmarshal(ext.Raw, &unstructBlob)
	if err != nil {
		return "", err
	}
	unstructuredOrig.Object = unstructBlob.(map[string]interface{})

	// Add placeholders to the object.
	currJSON := ext.Raw
	for _, opt := range tmplOpts {
		if opt.TemplateMatcher != nil && !opt.TemplateMatcher(unstructuredOrig.Object, resourceType) {
			continue
		}

		currJSON, err = addPlaceholder(opt, *gvk, currJSON)
		if err != nil {
			return "", err
		}
	}

	// Convert back to YAML.
	y, err := k8syaml.JSONToYAML(currJSON)
	if err != nil {
		return "", err
	}

	return string(y), nil
}

// TemplatizeK8sYAML takes a K8s YAML and templatizes the provided fields.
func TemplatizeK8sYAML(clientset *kubernetes.Clientset, inputYAML string, tmplOpts []*K8sTemplateOptions) (string, error) {
	// Read the YAML into K8s resources.
	yamlReader := strings.NewReader(inputYAML)

	decodedYAML := yaml.NewYAMLOrJSONDecoder(yamlReader, 4096)
	discoveryClient := clientset.Discovery()

	apiGroupResources, err := restmapper.GetAPIGroupResources(discoveryClient)
	if err != nil {
		return "", err
	}
	rm := restmapper.NewDiscoveryRESTMapper(apiGroupResources)

	templatedYAMLs := make([]string, 0)

	for { // Loop through all objects in the YAML and add placeholders.
		placeholderYAML, err := addPlaceholders(rm, decodedYAML, tmplOpts)
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return "", err
		}

		templatedYAMLs = append(templatedYAMLs, placeholderYAML)
	}
	// Concat the YAMLs.
	combinedYAML := ""
	for _, y := range templatedYAMLs {
		combinedYAML = ConcatYAMLs(combinedYAML, y)
	}

	// Replace all placeholders with their template values.
	replacedStrings := make([]string, 0)
	for _, opt := range tmplOpts {
		replacedStrings = append(replacedStrings, []string{opt.Placeholder, opt.TemplateValue}...)
	}

	r := strings.NewReplacer(replacedStrings...)

	return r.Replace(combinedYAML), nil
}
