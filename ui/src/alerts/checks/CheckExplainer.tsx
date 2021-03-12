import React from 'react'
import {ComponentSize, Gradients, Panel} from '@influxdata/clockface'

const CheckExplainer: React.FC = () => {
  return (
    <Panel
      gradient={Gradients.PolarExpress}
      border={true}
      style={{textAlign: 'inherit'}}
    >
      <Panel.Header>
        <h5>What is Check</h5>
      </Panel.Header>

      <Panel.Body size={ComponentSize.Small}>
        Check is blah blah blah
      </Panel.Body>
    </Panel>
  )
}

export default CheckExplainer
