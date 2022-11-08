import {ComponentSize, EmptyState} from '@influxdata/clockface'
import React, {FunctionComponent} from 'react'
import CreateConfigurationButton from './CreateConfigurationButton'

interface Props {
  searchTerm: string
}

const EmptyConfigurations: FunctionComponent<Props> = ({searchTerm}) => {
  if (searchTerm) {
    return (
      <EmptyState
        size={ComponentSize.Large}
        testID={'no-match-configurations-list'}
      >
        <EmptyState.Text>
          No Configurations match your search term <b>{searchTerm}</b>
        </EmptyState.Text>
      </EmptyState>
    )
  }

  return (
    <EmptyState size={ComponentSize.Large} testID={'empty-dashboards-list'}>
      <EmptyState.Text>
        Looks like you don't have any <b>Configurations</b>, why not create one?
      </EmptyState.Text>

      <CreateConfigurationButton />
    </EmptyState>
  )
}

export default EmptyConfigurations
