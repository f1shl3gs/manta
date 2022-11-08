import {Page} from '@influxdata/clockface'
import React, {FunctionComponent} from 'react'

const ConfigurationEditPage: FunctionComponent = () => {
  return (
    <Page>
      <Page.Header fullWidth={true}>head</Page.Header>

      <Page.Contents>contents</Page.Contents>
    </Page>
  )
}

export default ConfigurationEditPage
