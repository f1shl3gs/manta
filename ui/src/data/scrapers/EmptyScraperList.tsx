import React from 'react'
import {ComponentSize, EmptyState} from '@influxdata/clockface'

interface Props {
  searchTerm: string
}

const EmptyScraperList: React.FC<Props> = props => {
  const {searchTerm} = props

  if (searchTerm) {
    return (
      <EmptyState size={ComponentSize.Small} className={'alert-column--empty'}>
        <EmptyState.Text>
          No <b>Scrapers</b> match your search
        </EmptyState.Text>
      </EmptyState>
    )
  }

  return (
    <EmptyState size={ComponentSize.Small} className={'alert-column--empty'}>
      <EmptyState.Text>
        Looks like you have not create a <b>Scraper</b> yet
        <br />
        <br />
        You will need one to scrape metrics
      </EmptyState.Text>
    </EmptyState>
  )
}

export default EmptyScraperList
