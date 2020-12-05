import React from 'react';
import { Route, Switch, withRouter } from 'react-router-dom';

// styles
import './App.css';
import '@influxdata/clockface/dist/index.css';

// components
import { AppWrapper } from '@influxdata/clockface';
import Sidebar from 'layout/nav';
import Org from 'containers/organization/org';

const createOrg: React.FC = () => {
  return <div>create org</div>;
};

const App: React.FC = () => {
  return (
    <AppWrapper presentationMode={false} className="dark">
      <Sidebar />
      <Switch>
        <Route path="/orgs/new" component={createOrg} />
        <Route path="/orgs/:orgID" component={Org} />
      </Switch>
    </AppWrapper>
  );
};

export default withRouter(App);