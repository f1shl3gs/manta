// Libraries
import React, {FunctionComponent} from 'react'
import {useSelector} from 'react-redux'

// Components
import FilterList from 'src/shared/components/FilterList'
import CreateScrapeButton from 'src/scrapes/components/CreateScrapeButton'
import ScrapeCard from 'src/scrapes/components/ScrapeCard'

// Types
import {AppState} from 'src/types/stores'
import {Scrape} from 'src/types/scrape'
import {ResourceType} from 'src/types/resources'

// Selectors
import {getAll} from 'src/resources/selectors'
import EmptyResources from 'src/resources/components/EmptyResources'
import {getSortedResources} from 'src/shared/utils/sort'

interface Props {
  search: string
  sortOption: {
    key
    type
    direction
  }
}

const ScrapeList: FunctionComponent<Props> = ({search, sortOption}) => {
  const scrapes = useSelector((state: AppState) =>
    getAll<Scrape>(state, ResourceType.Scrapes)
  )

  return (
    <FilterList<Scrape>
      list={scrapes}
      search={search}
      searchKeys={['name', 'desc']}
    >
      {filtered => {
        if (filtered && filtered.length === 0) {
          return (
            <EmptyResources
              searchTerm={search}
              resource={ResourceType.Scrapes}
              createButton={<CreateScrapeButton />}
            />
          )
        }

        return (
          <div style={{height: '100%', display: 'grid'}}>
            {getSortedResources<Scrape>(
              filtered,
              sortOption.key,
              sortOption.type,
              sortOption.direction
            ).map(sc => (
              <ScrapeCard key={sc.id} scrape={sc} />
            ))}
          </div>
        )
      }}
    </FilterList>
  )
}

export default ScrapeList
