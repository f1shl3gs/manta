// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {Panel, Gradients} from '@influxdata/clockface'

const ScrapeExplainer: FunctionComponent = () => {
  return (
    <Panel gradient={Gradients.PolarExpress} border={true}>
      <Panel.Header>
        <h5>What is a Scrape?</h5>
      </Panel.Header>

      <Panel.Body>
        <p>Blah Blah Blah...</p>
      </Panel.Body>
    </Panel>
  )
}

export default ScrapeExplainer
