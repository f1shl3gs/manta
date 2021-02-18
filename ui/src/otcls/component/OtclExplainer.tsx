// Libraries
import React from 'react'

// Components
import {
  ComponentSize,
  EmptyState,
  Gradients,
  Panel,
  InfluxColors,
} from '@influxdata/clockface'
// @ts-ignore
import {TextAlignProperty} from 'csstype'

interface Props {
  hasOtcls?: boolean
  textAlign?: TextAlignProperty
  bodySize?: ComponentSize
}

const OtclExplainer: React.FC<Props> = props => {
  const {textAlign = 'inherit', hasOtcls = false, bodySize} = props

  return (
    <Panel
      gradient={Gradients.PolarExpress}
      border={true}
      style={{textAlign: textAlign}}
    >
      {!hasOtcls && (
        <EmptyState.Text style={{color: InfluxColors.Platinum, marginTop: 16}}>
          What is Otcl?
        </EmptyState.Text>
      )}
      {hasOtcls && (
        <Panel.Header>
          <h5>What is Otcl</h5>
        </Panel.Header>
      )}
      <Panel.Body size={bodySize}>
        Otcl is an agent written in Go for collection metrics, logs and tracing
        <br />
        <br />
        Here's a handy guide for <a href="/">Getting Started with Otcl</a>
      </Panel.Body>
    </Panel>
  )
}

export default OtclExplainer
