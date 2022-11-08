import {Gradients, Panel} from '@influxdata/clockface'
import React, {FunctionComponent} from 'react'

const ConfigurationExplainer: FunctionComponent = () => {
  return (
    <Panel gradient={Gradients.PolarExpress} border={true}>
      <Panel.Header>
        <h5>What is a Bucket?</h5>
      </Panel.Header>

      <Panel.Body>
        <p>Blah Blah Blah...</p>
      </Panel.Body>
    </Panel>
  )
}

export default ConfigurationExplainer
