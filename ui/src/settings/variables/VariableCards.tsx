// Libraries
import React from 'react'
import {Variable} from '../../types/Variable'
import {ResourceList, Sort} from '@influxdata/clockface'
import {getSortedResources} from '../../utils/sort'
import VariableCard from './VariableCard'
import {useVariables} from './useVariables'
import {SortTypes} from '../../types/sort'
import EmptyVariableList from './EmptyVariableList'

interface Props {
  list: Variable[]
  searchTerm: string
  sortKey: string
  sortType: SortTypes
  sortDirection: Sort
}

const VariableCards: React.FC<Props> = props => {
  const {list, searchTerm, sortKey, sortType, sortDirection} = props
  const {onNameUpdate, onDescUpdate, onDelete} = useVariables()

  const body = (filtered: Variable[]) => (
    <ResourceList.Body
      emptyState={<EmptyVariableList searchTerm={searchTerm} />}
    >
      {getSortedResources<Variable>(
        filtered,
        sortKey,
        sortType,
        sortDirection
      ).map(variable => (
        <VariableCard
          key={variable.id}
          variable={variable}
          onDelete={onDelete}
          onNameUpdate={onNameUpdate}
          onDescUpdate={onDescUpdate}
        />
      ))}
    </ResourceList.Body>
  )

  return <ResourceList>{body(list)}</ResourceList>
}

export default VariableCards
