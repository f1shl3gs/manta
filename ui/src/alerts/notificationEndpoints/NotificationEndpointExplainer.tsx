// Libraries
import React from 'react'

// Components
import {ComponentSize, Gradients, Panel} from '@influxdata/clockface'

const NotificationEndpointExplainer: React.FC = () => {
  return (
    <Panel
      gradient={Gradients.PolarExpress}
      border={true}
      style={{textAlign: 'inherit'}}
    >
      <Panel.Header>
        <h5>What is Notification Endpoint</h5>
      </Panel.Header>

      <Panel.Body size={ComponentSize.Small}>
        Notification endpoint is blah blah blah
      </Panel.Body>
    </Panel>
  )
}

export default NotificationEndpointExplainer
