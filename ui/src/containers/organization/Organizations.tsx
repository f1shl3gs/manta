// Library
import React from 'react'

// Components
import {
  RemoteDataState,
  SpinnerContainer,
  TechnoSpinner,
} from '@influxdata/clockface'
import {Redirect, Route, Switch} from 'react-router-dom'
import Org from './org'

// Hooks
import {useOrgs} from '../../shared/useOrgs'

const Organizations: React.FC = () => {
  const {orgs, loading} = useOrgs()

  return (
    <SpinnerContainer loading={loading} spinnerComponent={<TechnoSpinner />}>
      <Switch>
        <Route exact path={'/orgs'}>
          {loading === RemoteDataState.Done ? (
            <Redirect from="/" to={`/orgs/${orgs[0].id}`} />
          ) : null}
        </Route>

        <Route path={'/orgs/:orgID'} component={Org} />
      </Switch>
    </SpinnerContainer>
  )
}

export default Organizations
