import React, {FunctionComponent} from 'react'
import {Page} from '@influxdata/clockface'
import {Todo} from 'src/Todo'

const DashboardPage: FunctionComponent = () => {
  return (
    <Page>
      <Todo />
    </Page>
  )
}

export default DashboardPage
