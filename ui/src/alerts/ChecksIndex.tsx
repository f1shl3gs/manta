import React, {useState} from 'react'
import {useChecks} from './useChecks'
import TabbedPageHeader from '../shared/components/TabbedPageHeader'
import SearchWidget from '../shared/components/SearchWidget'
import {
  Button,
  ComponentColor,
  ComponentStatus,
  IconFont,
  Sort,
} from '@influxdata/clockface'
import FilterList from '../shared/components/FilterList'
import {Check} from '../types/Check'
import ResourceSortDropdown from '../shared/components/ResourceSortDropdown'
import {SortTypes} from '../types/Sort'

const ChecksIndex: React.FC = () => {
  const {checks, remoteDataState} = useChecks()
  const title = `Checks`
  const [search, setSearch] = useState('')

  const tooltipContents = (
    <>
      A <strong>Check</strong> is a periodic query that the system
      <br />
      performs against your time series data
      <br />
      that will generate a status
      <br />
      <br />
    </>
  )

  const [st, setSt] = useState('')

  const leftHeader = (
    <>
      <SearchWidget
        search={st}
        placeholder={'Filter Checks...'}
        onSearch={v => console.log('v', v)}
      />
      <ResourceSortDropdown
        sortKey={'updated'}
        sortType={SortTypes.Date}
        sortDirection={Sort.Ascending}
        onSelect={(sk, sd, st) => console.log(sk, sd, st)}
      />
    </>
  )

  const rightHeader = (
    <Button
      text={'Create Check'}
      icon={IconFont.Plus}
      color={ComponentColor.Primary}
      titleText={'Create a new Check'}
      status={ComponentStatus.Default}
    />
  )

  return (
    <>
      <TabbedPageHeader left={leftHeader} right={rightHeader} />

      <FilterList<Check>
        list={checks}
        search={''}
        searchKeys={['name', 'desc']}
      >
        {filtered => filtered.map(item => <div>{item.name}</div>)}
      </FilterList>
    </>
  )
}

export default ChecksIndex
