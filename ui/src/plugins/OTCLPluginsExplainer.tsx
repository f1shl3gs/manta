import React from 'react'
import {
  FontWeight,
  Heading,
  HeadingElement,
  InfluxColors,
  Panel,
} from '@influxdata/clockface'

const OTCLPluginsExplainer: React.FC = () => {
  return (
    <Panel backgroundColor={InfluxColors.Castle} style={{marginBottom: '8px'}}>
      <Panel.Header>
        <Heading element={HeadingElement.H4} weight={FontWeight.Regular}>
          Getting started with OTCL
        </Heading>
      </Panel.Header>

      <Panel.Body>
        <p>OTCL is blah blah blah</p>
      </Panel.Body>
    </Panel>
  )
}

export default OTCLPluginsExplainer
