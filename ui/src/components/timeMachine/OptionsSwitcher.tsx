import React from 'react';
import { ViewProperties } from 'types/Dashboard';
import LineOptions from './LineOptions';
import { ViewProvider } from './useView';

interface Props {
  view: ViewProperties
}

const OptionsSwitcher: React.FC<Props> = (props) => {
  const { view } = props;

  switch (view.type) {
    case 'gauge':
    case 'xy':
      return (
        <ViewProvider view={view}>
          <LineOptions />
        </ViewProvider>
      );
    default:
      return <div>Unknown</div>;
  }
};

export default OptionsSwitcher;
