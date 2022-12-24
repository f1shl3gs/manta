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

// Lazy load components
const ChecksIndex = lazy(() => import('src/checks/ChecksIndex'))
const CreateOrgOverlay = lazy(
  () => import('src/organizations/components/CreateOrgOverlay')
)
const DashboardImportOverlay = lazy(
  () => import('src/dashboards/components/DashboardImportOverlay')
)
const DashboardsIndex = lazy(
  () => import('src/dashboards/components/DashboardsIndex')
)
const DashboardPage = lazy(
  () => import('src/dashboards/components/DashboardPage')
)
const EditCheckOverlay = lazy(
  () => import('src/checks/components/EditCheckOverlay')
)
const Explore = lazy(() => import('src/explore/Explore'))
const Introduce = lazy(() => import('src/Introduce'))
const NewCheckOverlay = lazy(
  () => import('src/checks/components/NewCheckOverlay')
)
const SettingsPage = lazy(() => import('src/settings/SettingsIndex'))

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

                  <Route path="checks" element={<ChecksIndex />} />

                  <Route path={'checks/new'} element={<NewCheckOverlay />} />

                  <Route
                    path={'checks/:id/edit'}
                    element={<EditCheckOverlay />}
                  />

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
              <Route path="*" element={<NotFound />} />
            </Routes>
          </Suspense>
        </Organizations>
      </Authenticate>
    </AppWrapper>
  )
}

export default App
