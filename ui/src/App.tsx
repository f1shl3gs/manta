import "./wdyr"

import React from 'react';
import {
  Route,
  RouteComponentProps,
  Switch, withRouter
} from 'react-router-dom'

// styles
import './App.css';
import '@influxdata/clockface/dist/index.css';

// components
import {AppWrapper} from '@influxdata/clockface';
import Sidebar from 'layout/nav';
import Org from "containers/organization/org";

const createOrg: React.FC = props => {
  return <div>create org</div>
}

type Props = RouteComponentProps

const App: React.FC<Props> = props => {
  return (
    <AppWrapper
      presentationMode={false}
      className={'dark'}
    >
      <Sidebar/>
      <Switch>
        <Route path="/orgs/new" component={createOrg}/>
        <Route path="/orgs/:orgID" component={Org}/>
      </Switch>
    </AppWrapper>
  )
}

export default withRouter(App);
