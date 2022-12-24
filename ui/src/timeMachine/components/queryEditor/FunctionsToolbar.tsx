// Libraries
import React, {useMemo, useState} from 'react'

// Components
import {
  ComponentSize,
  DapperScrollbars,
  EmptyState,
} from '@influxdata/clockface'
import PromqlToolbarSearch from 'src/timeMachine/components/queryEditor/PromqlToolbarSearch'
import FunctionItem from 'src/timeMachine/components/queryEditor/FunctionItem'

// Constants
import {PROMQL_FUNCTIONS} from 'src/shared/constants/promqlFunctions'

const FunctionsToolbar: React.FC = () => {
  const [search, setSearch] = useState('')
  const filtered = useMemo(() => {
    return PROMQL_FUNCTIONS.filter(fn => {
      return fn.name.indexOf(search) >= 0
    })
  }, [search])

  let content: JSX.Element | JSX.Element[]
  if (filtered.length === 0) {
    content = (
      <EmptyState size={ComponentSize.ExtraSmall}>
        <EmptyState.Text>No functions match your search</EmptyState.Text>
      </EmptyState>
    )
  } else {
    content = filtered.map(fn => (
      <FunctionItem
        key={fn.name}
        fn={fn}
        testID={fn.name}
        onClickFn={cf => console.log('click', cf.name)}
      />
    ))
  }

  return (
    <>
      <PromqlToolbarSearch onSearch={setSearch} resourceName={'functions'} />

      <DapperScrollbars className={'flux-toolbar--scroll-area'}>
        <div className={'flux-toolbar--list'}>{content}</div>
      </DapperScrollbars>
    </>
  )
}

export default FunctionsToolbar
