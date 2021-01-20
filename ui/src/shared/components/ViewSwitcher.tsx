// Libraries
import React from 'react'

// Components
import {FromFluxResult, Plot} from '@influxdata/giraffe'
import XYPlot from '../../components/timeMachine/XYPlot'

// Types
import {ViewProperties} from 'types/Dashboard'
import GaugeChart from '../../components/timeMachine/GaugeChart'
import SingleStat from '../../components/timeMachine/SingleStat'
import LatestValueTransform from '../../components/timeMachine/LatestValueTransform'

interface Props {
  giraffeResult: Omit<FromFluxResult, 'schema'>
  properties: ViewProperties
}

const ViewSwitcher: React.FC<Props> = (props) => {
  const {
    properties,
    giraffeResult: {table, fluxGroupKeyUnion},
  } = props

  switch (properties.type) {
    case 'xy':
      return (
        <XYPlot table={table} groupKeyUnion={fluxGroupKeyUnion}>
          {(config) => <Plot config={config} />}
        </XYPlot>
      )

    case 'gauge':
      return (
        <LatestValueTransform table={table}>
          {(latestValue) => <GaugeChart value={latestValue} />}
        </LatestValueTransform>
      )
    case 'single-stat':
      return (
        <LatestValueTransform table={table}>
          {(latestValue) => (
            <SingleStat
              stat={latestValue}
              properties={properties}
              theme={'dark'}
            />
          )}
        </LatestValueTransform>
      )
    default:
      return <div>unknown</div>
  }
}

export default ViewSwitcher
