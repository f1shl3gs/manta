// Libraries
import React, {FunctionComponent} from 'react'
import {useSelector} from 'react-redux'

// Components
import FilterList from 'src/shared/components/FilterList'
import EmptyResources from 'src/resources/components/EmptyResources'
import CreateSecretButton from 'src/secrets/components/CreateSecretButton'
import SecretCard from 'src/secrets/components/SecretCard'

// Types
import {AppState} from 'src/types/stores'
import {Secret} from 'src/types/secrets'
import {ResourceType} from 'src/types/resources'

// Selectors
import {getAll} from 'src/resources/selectors'

// Utils
import {getSortedResources} from 'src/shared/utils/sort'

const SecretList: FunctionComponent = () => {
  const {secrets, searchTerm, sortOptions} = useSelector((state: AppState) => {
    const secrets = getAll<Secret>(state, ResourceType.Secrets)
    const {searchTerm, sortOptions} = state.resources[ResourceType.Secrets]

    return {
      secrets,
      searchTerm,
      sortOptions,
    }
  })

  return (
    <FilterList<Secret> list={secrets} search={searchTerm} searchKeys={['key']}>
      {filterted => {
        if (filterted && filterted.length === 0) {
          return (
            <EmptyResources
              resource={ResourceType.Secrets}
              createButton={<CreateSecretButton />}
            />
          )
        }

        return (
          <div>
            {getSortedResources(
              filterted,
              sortOptions.key,
              sortOptions.type,
              sortOptions.direction
            ).map(secret => (
              <SecretCard key={secret.key} secret={secret} />
            ))}
          </div>
        )
      }}
    </FilterList>
  )
}

export default SecretList
