// Libraries
import React, {FunctionComponent} from 'react'

// Components
import ErrorBoundary from 'src/shared/components/ErrorBoundary'

// Types
import {FromFluxResult} from '@influxdata/giraffe'
import {ViewProperties} from 'src/types/cells'
import Line from 'src/visualization/Line'

interface Props {
  result: FromFluxResult
  properties: ViewProperties
}

const View: FunctionComponent<Props> = ({result, properties}) => {
  const inner = () => {
    switch (properties.type) {
      case 'xy':
        return <Line properties={properties} result={result} />
      default:
        return <></>
    }
  }

  return <ErrorBoundary>{inner()}</ErrorBoundary>
}

export default View
