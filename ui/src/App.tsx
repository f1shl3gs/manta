// Libraries
import React, {FC, lazy, Suspense} from 'react'

// Components
import {AppWrapper} from '@influxdata/clockface'
import {usePresentationMode} from 'src/shared/usePresentationMode'
import {AuthenticationProvider} from 'src/shared/components/useAuthentication'
import {Route, Routes} from 'react-router-dom'
import Organizations from 'src/organizations/Organizations'
import {NotificationProvider} from 'src/shared/components/notifications/useNotification'
import Notifications from 'src/shared/components/notifications/Notifications'
import PageSpinner from 'src/shared/components/PageSpinner'
import Authentication from 'src/shared/components/Authentication'
import CreateOrgOverlay from 'src/organizations/CreateOrgOverlay'
import ToOrg from 'src/organizations/ToOrg'
// DataPage is just a simple tabed page, it's small enough and it can reduce re-render
import DataPage from 'src/data/DataPage'
import {AutoRefreshProvider} from 'src/shared/useAutoRefresh'
import {TimeRangeProvider} from 'src/shared/useTimeRange'
import DashboardPage from './dashboards/DashboardPage'

const Introduce = lazy(() => import('src/Introduce'))
const DashboardsPage = lazy(() => import('src/dashboards/DashboardsPage'))
const SettingsPage = lazy(() => import('src/settings/SettingsPage'))
const Explore = lazy(() => import('src/explore/Explore'))
const DashboardImportOverlay = lazy(() => import('src/dashboards/DashboardImportOverlay'))
const ExportOverlay = lazy(() => import('src/dashboards/ExportOverlay'))
const EditVEO = lazy(() => import('src/dashboards/EditVEO'))
const NewVEO = lazy(() => import('src/dashboards/NewVEO'))

const App: FC = () => {
  const {inPresentationMode} = usePresentationMode()

  return (
    <AppWrapper presentationMode={inPresentationMode}>
      <AuthenticationProvider>
        <NotificationProvider>
          <Authentication>
            <Notifications />

            <Organizations>
              <Suspense fallback={<PageSpinner />}>
                <TimeRangeProvider>
                  <AutoRefreshProvider>
                    <Routes>
                      <Route index element={<ToOrg />} />

                      <Route path="orgs">
                        <Route path="new" element={<CreateOrgOverlay />} />

                        <Route path=":orgId">
                          <Route index={true} element={<Introduce />} />

                          <Route path="data/*" element={<DataPage />} />

                          <Route path="explore" element={<Explore />} />

                          <Route
                            path="dashboards"
                            element={<DashboardsPage />}
                          />
                          <Route
                            path="dashboards/:dashboardId"
                            element={<DashboardPage />}
                          />
                          <Route
                            path="dashboards/import"
                            element={<DashboardImportOverlay />}
                          />
                          <Route
                            path="dashboards/:dashboardId/cells/new"
                            element={<NewVEO />}
                          />
                          <Route
                            path="dashboards/:dashboardId/cells/:cellID/edit"
                            element={<EditVEO />}
                          />
                          <Route
                            path={'dashboards/:dashboardId/export'}
                            element={<ExportOverlay />}
                          />

                          <Route path="settings/*" element={<SettingsPage />} />
                        </Route>
                      </Route>
                    </Routes>
                  </AutoRefreshProvider>
                </TimeRangeProvider>
              </Suspense>
            </Organizations>
          </Authentication>
        </NotificationProvider>
      </AuthenticationProvider>
    </AppWrapper>
  )
}

export default App
