// Libraries
import React from 'react';
import { Route, Switch, withRouter } from 'react-router-dom';

// Components
import { AppWrapper } from '@influxdata/clockface';
import Nav from 'layout/nav';
import Org from 'containers/organization/org';
import { usePresentationMode } from './shared/usePresentationMode';

// Styles
import './App.css';

const createOrg: React.FC = () => {
  return <div>create org</div>;
};

const App: React.FC = () => {
  const { inPresentationMode } = usePresentationMode();

  return (
    <AppWrapper presentationMode={inPresentationMode} className="dark">
      <Nav />
      <Switch>
        <Route path="/orgs/new" component={createOrg} />
        <Route path="/orgs/:orgID" component={Org} />
      </Switch>
    </AppWrapper>
  );
};

export default withRouter(App);
