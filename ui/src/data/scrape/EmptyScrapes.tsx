import React, {FunctionComponent} from 'react'
import {ComponentSize, EmptyState} from '@influxdata/clockface'
import CreateScrapeButton from './CreateScrapeButton'

interface Props {
  searchTerm: string
}

const EmptyScrapes: FunctionComponent<Props> = ({searchTerm}) => {
  if (searchTerm) {
    return (
      <EmptyState size={ComponentSize.Large} testID={'no-match-scrape-list'}>
        <EmptyState.Text>
          No Scrapes match your search term <b>{searchTerm}</b>
        </EmptyState.Text>
      </EmptyState>
    )
  }

  return (
    <EmptyState size={ComponentSize.Large} testID={'empty-scrape-list'}>
      <EmptyState.Text>
        Looks like youd don't have any <b>Scrape</b>, why not create some?
      </EmptyState.Text>

      <CreateScrapeButton />
    </EmptyState>
  )
}

export default EmptyScrapes
