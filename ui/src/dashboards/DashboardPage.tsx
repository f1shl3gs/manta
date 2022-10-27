import React, {FunctionComponent} from 'react'
import {Page} from '@influxdata/clockface'
import {Todo} from '../Todo'

const DashboardPage: FunctionComponent = () => {
  return (
    <Page>
      <Todo />
    </Page>
  )
}

export default DashboardPage
