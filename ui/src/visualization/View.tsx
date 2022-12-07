// Libraries
import React, {FunctionComponent} from 'react'

// Components
import ErrorBoundary from 'src/shared/components/ErrorBoundary'

// Types
import {FromFluxResult} from '@influxdata/giraffe'
import {ViewProperties} from 'src/types/cells'
import Line from 'src/visualization/Line'
import Gauge from 'src/visualization/Gauge'
import SingleStat from 'src/visualization/SingleStat'
import SingleStatPlusLine from 'src/visualization/SingleStatPlusLine/view'

interface Props {
  result: FromFluxResult
  properties: ViewProperties
}

const View: FunctionComponent<Props> = ({result, properties}) => {
  const inner = () => {
    switch (properties.type) {
      case 'xy':
        return <Line properties={properties} result={result} />

      case 'gauge':
        return <Gauge properties={properties} result={result} />

      case 'single-stat':
        return <SingleStat properties={properties} result={result} />

      case 'line-plus-single-stat':
        return <SingleStatPlusLine properties={properties} result={result} />

      default:
        return <>not implement</>
    }
  }

  return <ErrorBoundary>{inner()}</ErrorBoundary>
}

export default View
