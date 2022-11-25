import React from 'react'

import {Grid, Form, Dropdown} from '@influxdata/clockface'
import ColumnSelector from 'src/shared/components/ColumnSelector'

import TimeFormatSetting from 'src/visualization/TimeMachine/TimeFormatSetting'
import YAxisTitle from 'src/visualization/TimeMachine/YAxisTitle'
import YAxisBase from 'src/visualization/TimeMachine/YAxisBase'
import AxisAffixes from 'src/visualization/TimeMachine/AxisAffixes'
import useLineView from 'src/visualization/Line/useLineView'

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

const LineOptions: React.FC = () => {
  const {
    xColumn,
    onSetXColumn,
    yColumn,
    onSetYColumn,
    numericColumns,
    timeFormat,
    onSetTimeFormat,
    hoverDimension,
    onSetHoverDimension,
    axes: {
      y: {prefix = '', suffix = '', label = '', base = ''},
    },
    onSetYAxisLabel,
    onSetYAxisBase,
    onSetYAxisPrefix,
    onSetYAxisSuffix,
  } = useLineView()

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
