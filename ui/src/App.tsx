// Libraries
import React, {FC, lazy, Suspense} from 'react'

// Components
import {AppWrapper} from '@influxdata/clockface'
import {usePresentationMode} from 'shared/usePresentationMode'
import {AuthenticationProvider} from 'shared/components/useAuthentication'
import {Route, Routes} from 'react-router-dom'
import Organizations from 'organizations/Organizations'
import {NotificationProvider} from 'shared/components/notifications/useNotification'
import Notifications from 'shared/components/notifications/Notifications'
import PageSpinner from './shared/components/PageSpinner'

const Introduce = lazy(() => import('Introduce'))
const DashboardsPage = lazy(() => import('dashboards/DashboardsPage'))
const DashboardPage = lazy(() => import('dashboards/DashboardPage'))
const SettingsPage = lazy(() => import('settings/SettingsPage'))

const App: FC = () => {
  const {inPresentationMode} = usePresentationMode()

  return (
    <AppWrapper presentationMode={inPresentationMode}>
      <AuthenticationProvider>
        <NotificationProvider>
          <Notifications />

          <Organizations>
            <Suspense fallback={<PageSpinner />}>
              <Routes>
                <Route path="orgs/:orgId" element={<Introduce />} />
                <Route
                  path="orgs/:orgId/dashboards"
                  element={<DashboardsPage />}
                />
                <Route
                  path="orgs/:orgId/dashboards/:dashboardId"
                  element={<DashboardPage />}
                />
                <Route
                  path="orgs/:orgId/settings/*"
                  element={<SettingsPage />}
                />
              </Routes>
            </Suspense>
          </Organizations>
        </NotificationProvider>
      </AuthenticationProvider>
    </AppWrapper>
  )
}

export default App
