import React from 'react';
import { ViewProperties } from 'types/Dashboard';
import LineOptions from './LineOptions';

interface Props {
  view: ViewProperties
}

const OptionsSwitcher: React.FC<Props> = (props) => {
  const { view } = props;

  switch (view.type) {
    case 'gauge':
      return <div>todo</div>

    case 'xy':
      return (
        <LineOptions />
      );
    default:
      return <div>Unknown</div>;
  }
};

export default OptionsSwitcher;
