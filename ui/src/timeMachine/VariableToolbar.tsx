import React from 'react'
import {
  ComponentSize,
  DapperScrollbars,
  EmptyState,
} from '@influxdata/clockface'
import VariableItem from 'src/timeMachine/VariableItem'
import {Variable} from 'src/types/Variable'
import PromqlToolbarSearch from 'src/timeMachine/PromqlToolbarSearch'

const VariableToolbar: React.FC = () => {
  // todo: implement it
  const variables: Variable[] = [
    {
      id: 'foo',
      created: 'a',
      updated: 'b',
      name: 'foo',
      desc: 'foo desc',
      orgID: '0',
      type: 'query',
      value: 'label_values(up, instance)',
    },
    {
      id: 'bar',
      created: 'a',
      updated: 'b',
      name: 'bar',
      desc: 'bar desc',
      orgID: '0',
      type: 'static',
      value: 'apple,banana',
    },
  ]

  let content: JSX.Element | JSX.Element[]

  if (variables.length !== 0) {
    content = variables.map(variable => (
      <VariableItem
        key={variable.id}
        variable={variable}
        testID={`${variable.name}`}
      />
    ))
  } else {
    content = (
      <EmptyState size={ComponentSize.ExtraSmall}>
        <EmptyState.Text>No variables match your search</EmptyState.Text>
      </EmptyState>
    )
  }

  return (
    <>
      <PromqlToolbarSearch
        resourceName={'variables'}
        onSearch={v => console.log('onsearch', v)}
      />

      <DapperScrollbars className={'flux-toolbar--scroll-area'}>
        <div className={'flux-toolbar--list'}>{content}</div>
      </DapperScrollbars>
    </>
  )
}

export default VariableToolbar
