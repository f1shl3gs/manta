import React from 'react';

import {
  Config,
  Table,
} from '@influxdata/giraffe';
import { XYViewProperties } from 'types/Dashboard';

interface Props {
  children: (config: Config) => JSX.Element
  // timeRange?: TimeRange | null
  table: Table
  viewProperties: XYViewProperties
}

const XYPlot: React.FC<Props> = props => {
  const {
    children,
    table,
    viewProperties: {
      timeFormat,
      xColumn = '_time',
      yColumn = '_value',
      axes: {
        x: {
          label: xAxisLabel,
          prefix: xTickPrefix
        },
        y: {
          label: yAxisLabel,
          prefix: yTickPrefix
        }
      }
    }
  } = props;

  const groupKey = [ 'result']

  const config: Config = {
    table,
    xAxisLabel,
    yAxisLabel,
    layers: [
      {
        type: 'line',
        x: xColumn,
        y: yColumn,
        fill: groupKey,
      }
    ]
  }

  return children(config)
};

export default XYPlot;