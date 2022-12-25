// Libraries
import React, {FunctionComponent, useCallback} from 'react'

// Components
import {Grid, Form, Dropdown} from '@influxdata/clockface'
import ColumnSelector from 'src/shared/components/ColumnSelector'
import TimeFormatSetting from 'src/timeMachine/components/TimeFormatSetting'
import YAxisTitle from 'src/timeMachine/components/YAxisTitle'
import YAxisBase from 'src/timeMachine/components/YAxisBase'
import AxisAffixes from 'src/timeMachine/components/AxisAffixes'
import {VisualizationOptionProps} from 'src/visualization'

// Cells
import {XYViewProperties} from 'src/types/cells'

const dimensions = [
  {
    key: 'auto',
    text: 'Auto',
  },
  {
    key: 'x',
    text: 'X',
  },
  {
    key: 'y',
    text: 'Y',
  },
  {
    key: 'xy',
    text: 'XY',
  },
]

interface Props extends VisualizationOptionProps {
  viewProperties: XYViewProperties
}

const LineOptions: FunctionComponent<Props> = ({viewProperties, update}) => {
  const {
    xColumn = '_time',
    yColumn = '_value',
    timeFormat,
    hoverDimension,
    axes: {
      y: {prefix = '', suffix = '', label = '', base = ''},
    },
  } = viewProperties
  const numericColumns = [xColumn, yColumn]

  const onSetXColumn = useCallback(
    (x: string) => {
      update({
        ...viewProperties,
        xColumn: x,
      })
    },
    [viewProperties, update]
  )

  const onSetYColumn = useCallback(
    (y: string) => {
      update({
        ...viewProperties,
        yColumn: y,
      })
    },
    [viewProperties, update]
  )

  const onSetTimeFormat = useCallback(
    (timeFormat: string) => {
      update({
        ...viewProperties,
        timeFormat,
      })
    },
    [viewProperties, update]
  )

  const onSetHoverDimension = useCallback(
    (hoverDimension: 'x' | 'y' | 'xy' | 'auto') => {
      update({
        ...viewProperties,
        hoverDimension,
      })
    },
    [viewProperties, update]
  )

  const updateYAxis = useCallback(
    (upd: {[key: string]: string}) => {
      update({
        ...viewProperties,
        axes: {
          x: viewProperties.axes.x,
          y: {
            ...viewProperties.axes.y,
            ...upd,
          },
        },
      })
    },
    [update, viewProperties]
  )

  const onSetYAxisLabel = useCallback(
    (label: string) => {
      updateYAxis({label})
    },
    [updateYAxis]
  )

  const onSetYAxisBase = useCallback(
    (base: string) => {
      updateYAxis({base})
    },
    [updateYAxis]
  )

  const onSetYAxisPrefix = useCallback(
    (prefix: string) => {
      updateYAxis({prefix})
    },
    [updateYAxis]
  )

  const onSetYAxisSuffix = useCallback(
    (suffix: string) => {
      updateYAxis({suffix})
    },
    [updateYAxis]
  )

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

        <Form.Element label={'Time Format'}>
          <TimeFormatSetting
            timeFormat={timeFormat}
            onTimeFormatChange={onSetTimeFormat}
          />
        </Form.Element>

        <h5 className={'view-options--header'}>Options</h5>
      </Grid.Column>

      <Grid.Column>
        <br />
        <Form.Element label={'Hover Dimension'}>
          <Dropdown
            button={(active, onClick) => (
              <Dropdown.Button active={active} onClick={onClick}>
                {hoverDimension}
              </Dropdown.Button>
            )}
            menu={onCollapse => (
              <Dropdown.Menu onCollapse={onCollapse}>
                {dimensions.map(item => (
                  <Dropdown.Item
                    id={item.key}
                    key={item.key}
                    value={item.key}
                    onClick={onSetHoverDimension}
                    selected={hoverDimension === item.key}
                  >
                    {item.text}
                  </Dropdown.Item>
                ))}
              </Dropdown.Menu>
            )}
          />
        </Form.Element>
      </Grid.Column>

      <Grid.Column>
        <h5 className={'view-options--header'}>Y Axis</h5>
      </Grid.Column>
      <YAxisTitle label={label} onUpdateYAxisLabel={onSetYAxisLabel} />
      <YAxisBase base={base} onSetYAxisBase={onSetYAxisBase} />
      <AxisAffixes
        prefix={prefix}
        suffix={suffix}
        axisName={'y'}
        onSetAxisPrefix={onSetYAxisPrefix}
        onSetAxisSuffix={onSetYAxisSuffix}
      />
    </>
  )
}

export default LineOptions
