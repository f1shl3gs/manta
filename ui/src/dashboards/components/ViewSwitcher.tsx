import React from 'react';
import { FromFluxResult, Plot } from '@influxdata/giraffe';
import { ViewProperties } from 'types/Dashboard';
import XYPlot from './XYPlot';

interface Props {
  giraffeResult: Omit<FromFluxResult, 'schema'>
  properties: ViewProperties
}

const ViewSwitcher: React.FC<Props> = props => {
  const {
    properties,
    giraffeResult: {
      table
    }
  } = props;

  switch (properties.type) {
    case 'xy':
      return (
        <XYPlot
          table={table}
          viewProperties={properties}>
          {config => <Plot config={config} />}
        </XYPlot>
      );
    default:
      return <div>
        unknown
      </div>;
  }
};

export default ViewSwitcher;
