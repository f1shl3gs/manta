import React, {FunctionComponent, lazy, useMemo} from 'react'
import {
  Orientation,
  Page,
  PageContents,
  PageHeader,
  PageTitle,
  Tabs,
  TabsContainer,
} from '@influxdata/clockface'
import {Route, Routes, useNavigate} from 'react-router-dom'

const Members = lazy(() => import('settings/Members'))
const Secrets = lazy(() => import('settings/Secrets'))

const SettingsPage: FunctionComponent = () => {
  const navigate = useNavigate()
  const pathname = window.location.pathname
  const selected = useMemo(
    () => pathname.split('/').pop() as string,
    [pathname]
  )

  return (
    <Page titleTag={`Settings | ${selected}`}>
      <PageHeader fullWidth={false}>
        <PageTitle title="Settings" />
      </PageHeader>

      <PageContents>
        <TabsContainer orientation={Orientation.Horizontal}>
          <Tabs>
            {['members', 'secrets'].map(key => (
              <Tabs.Tab
                key={key}
                active={selected === key}
                id={key}
                text={key}
                onClick={() => {
                  if (selected === key) {
                    return
                  }

                  navigate(window.location.pathname.replace(selected, key))
                }}
              />
            ))}
          </Tabs>

          <Tabs.TabContents>
            <Routes>
              <Route path={'members'} element={<Members />} />
              <Route path={'secrets'} element={<Secrets />} />
            </Routes>
          </Tabs.TabContents>
        </TabsContainer>
      </PageContents>
    </Page>
  )
}

export default SettingsPage
