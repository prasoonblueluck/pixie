import clsx from 'clsx';
import * as React from 'react';

import { createStyles, makeStyles, Theme } from '@material-ui/core/styles';

interface LazyPanelProps {
  show: boolean;
  className?: string;
  children: React.ReactNode;
}

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    panel: {
      '&:not(.visible)': {
        display: 'none',
      },
    },
  }));

// LazyPanel is a component that renders the content lazily.
const LazyPanel = (props: LazyPanelProps) => {
  const { show, className, children } = props;
  const [rendered, setRendered] = React.useState(false);
  const classes = useStyles();

  React.useEffect(() => {
    setTimeout(() => window.dispatchEvent(new Event('resize')), 0);
  }, [show]);

  if (!show && !rendered) {
    return null;
  }
  if (show && !rendered) {
    setRendered(true);
  }

  return (
    <div className={clsx(className, classes.panel, show && 'visible')}>
      {children}
    </div>
  );
};

export default LazyPanel;
