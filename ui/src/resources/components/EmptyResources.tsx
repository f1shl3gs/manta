// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {ComponentSize, EmptyState} from '@influxdata/clockface'

// Types
import {ResourceType} from 'src/types/resources'

interface Props {
  resource: ResourceType
  createButton: JSX.Element

  searchTerm?: string
}

const EmptyResources: FunctionComponent<Props> = ({
  resource,
  searchTerm,
  createButton,
}) => {
  if (searchTerm && searchTerm !== '') {
    return (
      <EmptyState size={ComponentSize.Large} testID={'no-match-checks-list'}>
        <EmptyState.Text>
          No {resource} match your search term <b>{searchTerm}</b>
        </EmptyState.Text>
      </EmptyState>
    )
  }

  return (
    <EmptyState size={ComponentSize.Large} testID={'empty-checks-list'}>
      <EmptyState.Text>
        Looks like you don't have any <b>{resource}</b>, why not create one?
      </EmptyState.Text>

      {createButton}
    </EmptyState>
  )
}

export default EmptyResources
