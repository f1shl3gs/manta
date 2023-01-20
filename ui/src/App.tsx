// Libraries
import React, {FunctionComponent, lazy, Suspense} from 'react'

// Components
import {AppWrapper} from '@influxdata/clockface'
import {Route, Routes} from 'react-router-dom'
import Organizations from 'src/organizations/components/Organizations'
import Notifications from 'src/shared/components/notifications/Notifications'
import PageSpinner from 'src/shared/components/PageSpinner'
import ToOrg from 'src/organizations/components/ToOrg'
// DataPage is just a simple tabed page, it's small enough and it can reduce re-render
import DataPage from 'src/data/DataPage'
import Authenticate from 'src/me/components/Authenticate'
import NotFound from 'src/shared/components/NotFound'

// Hooks
import {useSelector} from 'react-redux'

// Selectors
import {getPresentationMode} from 'src/shared/selectors/app'
import AlertsIndex from './alerts/AlertsIndex'

// Lazy load components
const ChecksIndex = lazy(() => import('src/checks/components/ChecksIndex'))
const CreateOrgOverlay = lazy(
  () => import('src/organizations/components/CreateOrgOverlay')
)
const DashboardsIndex = lazy(
  () => import('src/dashboards/components/DashboardsIndex')
)
const Explore = lazy(() => import('src/explore/Explore'))
const Introduce = lazy(() => import('src/Introduce'))
const SettingsPage = lazy(() => import('src/settings/SettingsIndex'))
const DashboardPage = lazy(
  () => import('src/dashboards/components/DashboardPage')
)
const DashboardImportOverlay = lazy(
  () => import('src/dashboards/components/DashboardImportOverlay')
)
const DashboardExportOverlay = lazy(
  () => import('src/dashboards/components/ExportOverlay')
)

const App: FunctionComponent = () => {
  const presentationMode = useSelector(getPresentationMode)

  return (
    <AppWrapper presentationMode={presentationMode}>
      <Authenticate>
        <Notifications />

        <Organizations>
          <Suspense fallback={<PageSpinner />}>
            <Routes>
              <Route index element={<ToOrg />} />

              <Route path="orgs">
                <Route path="new" element={<CreateOrgOverlay />} />

                <Route path=":orgID">
                  <Route index={true} element={<Introduce />} />

                  <Route path="alerts/*" element={<AlertsIndex />} />

                  <Route path="checks/*" element={<ChecksIndex />} />

                  <Route path="data/*" element={<DataPage />} />

                  <Route path="explore" element={<Explore />} />

                  <Route path="dashboards/*" element={<DashboardsIndex />} />
                  <Route
                    path="dashboards/import"
                    element={<DashboardImportOverlay />}
                  />

                  <Route
                    path="dashboards/:dashboardID/export"
                    element={<DashboardExportOverlay />}
                  />
                  <Route
                    path="dashboards/:dashboardID/*"
                    element={<DashboardPage />}
                  />

                  <Route path="settings/*" element={<SettingsPage />} />
                </Route>
              </Route>
              <Route path="*" element={<NotFound />} />
            </Routes>
          </Suspense>
        </Organizations>
      </Authenticate>
    </AppWrapper>
  )
}

export default App
