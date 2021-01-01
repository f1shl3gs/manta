// Libraries
import React from 'react';

// Components
import { FromFluxResult, Plot } from '@influxdata/giraffe';
import XYPlot from '../../components/timeMachine/XYPlot';

// Types
import { ViewProperties } from 'types/Dashboard';

interface Props {
  giraffeResult: Omit<FromFluxResult, 'schema'>
  properties: ViewProperties
}

const ViewSwitcher: React.FC<Props> = props => {
  const {
    properties,
    giraffeResult: {
      table,
      fluxGroupKeyUnion
    }
  } = props;

  switch (properties.type) {
    case 'xy':
      return (
        <XYPlot
          table={table}
          groupKeyUnion={fluxGroupKeyUnion}
        >
          {config => <Plot config={config} />}
        </XYPlot>
      );
      
    case 'gauge':
      return (
        <div>todo: Gauge</div>
      );

    default:
      return <div>
        unknown
      </div>;
  }
};

export default ViewSwitcher;
