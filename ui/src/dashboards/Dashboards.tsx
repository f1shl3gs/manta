// Libraries
import React, {useState} from 'react'
import {useHistory} from 'react-router-dom'

// Components
import {
  Button,
  ComponentColor,
  IconFont,
  Page,
  Sort,
} from '@influxdata/clockface'
import DashboardCards from './DashboardCards'
import SearchWidget from '../shared/components/SearchWidget'

// Hooks
import {DashboardsProvider} from './useDashboards'
import {useFetch} from 'shared/useFetch'
import {useOrgID} from 'shared/useOrg'

// Types
import {Dashboard} from 'types/Dashboard'
import ResourceSortDropdown from '../shared/components/ResourceSortDropdown'
import {SortKey, SortTypes} from '../types/sort'
import FilterList from '../shared/components/FilterList'

const useCreateDash = () => {
  const orgID = useOrgID()

  const {post} = useFetch<Dashboard>(`/api/v1/dashboards?orgID=${orgID}`, {
    body: {
      orgID,
    },
  })

  return post
}

const Dashboards: React.FC = () => {
  const create = useCreateDash()
  const history = useHistory()
  const orgID = useOrgID()
  const [searchTerm, setSearchTerm] = useState('')
  const [sortOption, setSortOption] = useState({
    key: 'updated' as SortKey,
    type: SortTypes.Date,
    direction: Sort.Descending,
  })

  return (
    <Page titleTag={'Dashboards'}>
      <Page.Header fullWidth={false}>
        <Page.Title title={'Dashboards'} />
        {/* rateLimitAlert? */}
      </Page.Header>

      <Page.ControlBar fullWidth={false}>
        <Page.ControlBarLeft>
          <SearchWidget
            search={searchTerm}
            placeholder={'Filter dashboards...'}
            onSearch={setSearchTerm}
          />
          <ResourceSortDropdown
            sortKey={sortOption.key}
            sortType={sortOption.type}
            sortDirection={sortOption.direction}
            onSelect={(sk, sd, st) => {
              setSortOption({
                key: sk,
                type: st,
                direction: sd,
              })
            }}
          />
        </Page.ControlBarLeft>
        <Page.ControlBarRight>
          <Button
            text={'Add'}
            icon={IconFont.Plus}
            color={ComponentColor.Primary}
            onClick={() => {
              create()
                .then(resp => {
                  const path = `/orgs/${orgID}/dashboards/${resp.id}`
                  console.log('path', path)
                  history.push(path)
                })
                .catch(err => {
                  console.log('create dashboard failed', err)
                })
            }}
          />
        </Page.ControlBarRight>
      </Page.ControlBar>

      <Page.Contents
        className="dashboards-index__page-contents"
        fullWidth={false}
        scrollable={true}
      >
        <DashboardsProvider>
          <DashboardCards
            searchTerm={searchTerm}
            sortKey={sortOption.key}
            sortType={sortOption.type}
            sortDirection={sortOption.direction}
          />
        </DashboardsProvider>
      </Page.Contents>
    </Page>
  )
}

export default Dashboards
