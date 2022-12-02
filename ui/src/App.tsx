// Libraries
import React, {FC, lazy, Suspense} from 'react'

// Components
import {AppWrapper} from '@influxdata/clockface'
import {AuthenticationProvider} from 'src/shared/components/useAuthentication'
import {Route, Routes} from 'react-router-dom'
import Organizations from 'src/organizations/Organizations'
import Notifications from 'src/shared/components/notifications/Notifications'
import PageSpinner from 'src/shared/components/PageSpinner'
import Authentication from 'src/shared/components/Authentication'
import CreateOrgOverlay from 'src/organizations/CreateOrgOverlay'
import ToOrg from 'src/organizations/ToOrg'
// DataPage is just a simple tabed page, it's small enough and it can reduce re-render
import DataPage from 'src/data/DataPage'
import DashboardPage from 'src/dashboards/components/DashboardPage'
import {getPresentationMode} from 'src/shared/selectors/app'

// Hooks
import {useSelector} from 'react-redux'

// Lazy load components
const Introduce = lazy(() => import('src/Introduce'))
const DashboardsPage = lazy(
  () => import('src/dashboards/components/DashboardsIndex')
)
const SettingsPage = lazy(() => import('src/settings/SettingsIndex'))
const Explore = lazy(() => import('src/explore/Explore'))
const DashboardImportOverlay = lazy(
  () => import('src/dashboards/components/DashboardImportOverlay')
)
const ExportOverlay = lazy(
  () => import('src/dashboards/components/ExportOverlay')
)
const EditVEO = lazy(() => import('src/dashboards/components/EditVEO'))
const NewVEO = lazy(() => import('src/dashboards/components/NewVEO'))

const App: FC = () => {
  const presentationMode = useSelector(getPresentationMode)

  return (
    <AppWrapper presentationMode={presentationMode}>
      <AuthenticationProvider>
        <Authentication>
          <Notifications />

          <Organizations>
            <Suspense fallback={<PageSpinner />}>
              <Routes>
                <Route index element={<ToOrg />} />

                <Route path="orgs">
                  <Route path="new" element={<CreateOrgOverlay />} />

                  <Route path=":orgID">
                    <Route index={true} element={<Introduce />} />

                    <Route path="data/*" element={<DataPage />} />

                    <Route path="explore" element={<Explore />} />

                    <Route path="dashboards" element={<DashboardsPage />} />
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
            </Suspense>
          </Organizations>
        </Authentication>
      </AuthenticationProvider>
    </AppWrapper>
  )
}

export default App
