// Libraries
import React, {FunctionComponent} from 'react'
import {useSelector} from 'react-redux'

// Components
import FilterList from 'src/shared/components/FilterList'
import EmptyResources from 'src/resources/components/EmptyResources'
import CreateCheckButton from 'src/checks/components/CreateCheckButton'
import CheckCard from 'src/checks/components/CheckCard'

// Types
import {AppState} from 'src/types/stores'
import {Check, CheckSortParams} from 'src/types/checks'
import {ResourceType} from 'src/types/resources'

// Selectors
import {getAll} from 'src/resources/selectors'

// Utils
import {getSortedResources} from 'src/shared/utils/sort'

interface Props {
  searchTerm: string
  sortOption: CheckSortParams
}

const CheckCards: FunctionComponent<Props> = ({searchTerm, sortOption}) => {
  const checks = useSelector((state: AppState) =>
    getAll<Check>(state, ResourceType.Checks)
  )

  return (
    <FilterList<Check>
      list={checks}
      search={searchTerm}
      searchKeys={['name', 'desc']}
    >
      {filtered => {
        if (filtered && filtered.length === 0) {
          return (
            <EmptyResources
              searchTerm={searchTerm}
              resource={ResourceType.Checks}
              createButton={<CreateCheckButton />}
            />
          )
        }

        return (
          <div style={{height: '100%', display: 'grid'}}>
            <div>
              {getSortedResources<Check>(
                filtered,
                sortOption.key,
                sortOption.type,
                sortOption.direction
              ).map(check => (
                <CheckCard key={check.id} check={check} />
              ))}
            </div>
          </div>
        )
      }}
    </FilterList>
  )
}

export default CheckCards
