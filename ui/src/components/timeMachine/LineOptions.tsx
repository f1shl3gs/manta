import React from 'react';

import { LinePosition } from '@influxdata/giraffe';
import { Grid } from '@influxdata/clockface';
import ColumnSelector from '../ColumnSelector';

import { useLineView } from './useView';

import { Axes, ViewType, XYGeom } from 'types/Dashboard';

interface X {
  xColumn?: string
  onSetXColumn: (c: string) => void
  yColumn?: string
  onSetYColumn: (c: string) => void
  // todo: figure out what it is
  numericColumns: string[]
}

interface Props extends X {
  type: ViewType
  axes: Axes
  geom?: XYGeom
  shadeBelow?: boolean
  hoverDimension?: 'auto' | 'x' | 'y' | 'xy'
  position: LinePosition
}

const LineOptions: React.FC = () => {
  const {
    xColumn,
    onSetXColumn,
    yColumn,
    onSetYColumn,
    numericColumns
  } = useLineView();

  return (
    <>
      <Grid.Column>
        <h4 className={'view-options--header'}>Customize Line Graph</h4>
        <h5 className={'view-options--header'}>Data</h5>
        <ColumnSelector
          selectedColumn={xColumn}
          onSelectColumn={onSetXColumn}
          availableColumns={numericColumns}
          axisName={'x'}
        />
        <ColumnSelector
          selectedColumn={yColumn}
          onSelectColumn={onSetYColumn}
          availableColumns={numericColumns}
          axisName={'y'}
        />
      </Grid.Column>
    </>
  );
};

export default LineOptions;
