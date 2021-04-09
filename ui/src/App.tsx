// Libraries
import React from 'react'
import {Redirect, Route, Switch, withRouter} from 'react-router-dom'

// Components
import {AppWrapper} from '@influxdata/clockface'
import {usePresentationMode} from './shared/usePresentationMode'
import Organizations from './containers/organization/Organizations'
import Signin from './components/Signin'
import Notifications from './components/notifications/Notifications'
import Logout from './Logout'

// Styles

const App: React.FC = () => {
  const {inPresentationMode} = usePresentationMode()

  return (
    <AppWrapper presentationMode={inPresentationMode} className="dark">
      <Notifications />
      <Switch>
        <Redirect exact from={'/'} to={'/orgs'} />

        <Route path={'/orgs'} component={Organizations} />
        <Route path={'/signin'} component={Signin} />
        <Route path={'/logout'} component={Logout} />
      </Switch>
    </AppWrapper>
  )
}

export default withRouter(App)
