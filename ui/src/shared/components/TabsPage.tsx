import React, {FunctionComponent, useMemo} from 'react'
import {Route, Routes, useNavigate} from 'react-router-dom'
import {
  Orientation,
  Page,
  PageContents,
  PageHeader,
  PageTitle,
  Tabs,
  TabsContainer,
} from '@influxdata/clockface'

interface Tab {
  name: string
  element: JSX.Element
}

interface Props {
  title: string
  tabs: Tab[]
}

const TabsPage: FunctionComponent<Props> = ({title, tabs}) => {
  const navigate = useNavigate()
  const pathname = window.location.pathname
  const selected = useMemo(() => pathname.split('/')[4]!, [pathname])

  return (
    <Page titleTag={`${title} | ${selected}`}>
      <PageHeader fullWidth={false}>
        <PageTitle title={title} />
      </PageHeader>

      <PageContents>
        <TabsContainer orientation={Orientation.Horizontal}>
          <Tabs>
            {tabs.map(({name}) => (
              <Tabs.Tab
                key={name}
                active={selected === name}
                id={name}
                text={name}
                testID={`tab-${name}`}
                onClick={() => {
                  if (selected === name) {
                    return
                  }

                  navigate(window.location.pathname.replace(selected, name))
                }}
              />
            ))}
          </Tabs>

          <Tabs.TabContents>
            <Routes>
              {tabs.map(({name, element}) => (
                <Route key={name} path={`${name}`} element={element} />
              ))}
            </Routes>
          </Tabs.TabContents>
        </TabsContainer>
      </PageContents>
    </Page>
  )
}

export default TabsPage
