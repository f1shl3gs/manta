// Libraries
import React from 'react'
import {
  ComponentSize,
  EmptyState,
  Gradients,
  InfluxColors,
  Panel,
} from '@influxdata/clockface'

// Components

const ScraperExplainer: React.FC = () => {
  const textAlign = 'inherit'
  const hasResource = false
  const bodySize = ComponentSize.Medium

  return (
    <Panel
      gradient={Gradients.PolarExpress}
      border={true}
      style={{textAlign: textAlign}}
    >
      {!hasResource && (
        <EmptyState.Text style={{color: InfluxColors.Platinum, marginTop: 16}}>
          What is Otcl?
        </EmptyState.Text>
      )}
      {hasResource && (
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

export default ScraperExplainer
