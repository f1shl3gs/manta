// Libraries
import React from 'react'
import {Route, Switch, useHistory, withRouter} from 'react-router-dom'

// Components
import {AppWrapper} from '@influxdata/clockface'
import Nav from 'layout/nav'
import Org from 'containers/organization/org'
import {usePresentationMode} from './shared/usePresentationMode'

// Styles
import './App.css'
import {Provider} from 'use-http'

const createOrg: React.FC = () => {
  return <div>create org</div>
}

const App: React.FC = () => {
  const {inPresentationMode} = usePresentationMode()
  const history = useHistory()

  const options = {
    interceptors: {
      // @ts-ignore
      response: async ({response}) => {
        console.log('interceptors')

        if (response === undefined) {
          return undefined
        }

        if (response.status === 401) {
          history.push('/signin')
          return
        }

        return response
      },
    },
  }

  return (
    <AppWrapper presentationMode={inPresentationMode} className="dark">
      <Nav />

      <Switch>
        <Route path="/orgs/new" component={createOrg} />
        <Route path="/orgs/:orgID" component={Org} />
      </Switch>
    </AppWrapper>
  )
}

export default withRouter(App)
