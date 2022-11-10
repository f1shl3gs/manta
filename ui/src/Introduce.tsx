import {Page} from '@influxdata/clockface'
import React, {FunctionComponent} from 'react'

const Introduce: FunctionComponent = () => {
  return (
    <Page titleTag={'Introduce'}>
      <Page.Header fullWidth={true}>
        <Page.Title title={'Getting started'} />
      </Page.Header>

      <Page.Contents>Todo</Page.Contents>
    </Page>
  )
}

export default Introduce
