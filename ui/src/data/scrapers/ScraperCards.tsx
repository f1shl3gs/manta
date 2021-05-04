import React from 'react'
import {SortKey, SortTypes} from '../../types/sort'
import {ResourceList, Sort} from '@influxdata/clockface'
import {Scraper} from '../../types/scrapers'
import {getSortedResources} from '../../utils/sort'
import ScraperCard from './ScraperCard'
import EmptyScraperList from './EmptyScraperList'

interface Props<T> {
  list: T[]

  searchTerm: string
  sortKey: SortKey
  sortType: SortTypes
  sortDirection: Sort
}

const ScraperCards: React.FC<Props<Scraper>> = props => {
  const {list, searchTerm, sortKey, sortType, sortDirection} = props

  const body = (filtered: Scraper[]) => (
    <ResourceList.Body
      emptyState={<EmptyScraperList searchTerm={searchTerm} />}
    >
      {getSortedResources<Scraper>(
        filtered,
        sortKey,
        sortType,
        sortDirection
      ).map(scraper => (
        <ScraperCard key={scraper.id} scraper={scraper} />
      ))}
    </ResourceList.Body>
  )

  return <ResourceList>{body(list)}</ResourceList>
}

export default ScraperCards
