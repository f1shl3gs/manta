// Libraries
import React, {FC, lazy, Suspense} from 'react'

// Components
import {AppWrapper} from '@influxdata/clockface'
import {Route, Routes} from 'react-router-dom'
import Organizations from 'src/organizations/components/Organizations'
import Notifications from 'src/shared/components/notifications/Notifications'
import PageSpinner from 'src/shared/components/PageSpinner'
import CreateOrgOverlay from 'src/organizations/components/CreateOrgOverlay'
import ToOrg from 'src/organizations/components/ToOrg'
// DataPage is just a simple tabed page, it's small enough and it can reduce re-render
import DataPage from 'src/data/DataPage'
import {getPresentationMode} from 'src/shared/selectors/app'
import DashboardImportOverlay from './dashboards/components/DashboardImportOverlay'
import Authenticate from 'src/me/components/Authenticate'

// Hooks
import {useSelector} from 'react-redux'

// Lazy load components
const Introduce = lazy(() => import('src/Introduce'))
const DashboardsIndex = lazy(
  () => import('src/dashboards/components/DashboardsIndex')
)
const DashboardPage = lazy(
  () => import('src/dashboards/components/DashboardPage')
)
const SettingsPage = lazy(() => import('src/settings/SettingsIndex'))
const Explore = lazy(() => import('src/explore/Explore'))

const App: FC = () => {
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

                  <Route path="data/*" element={<DataPage />} />

                  <Route path="explore" element={<Explore />} />

                  <Route path="dashboards/*" element={<DashboardsIndex />} />

                  <Route
                    path="dashboards/import"
                    element={<DashboardImportOverlay />}
                  />

                  <Route
                    path="dashboards/:dashboardID/*"
                    element={<DashboardPage />}
                  />

                  <Route path="settings/*" element={<SettingsPage />} />
                </Route>
              </Route>
            </Routes>
          </Suspense>
        </Organizations>
      </Authenticate>
    </AppWrapper>
  )
}

export default App
