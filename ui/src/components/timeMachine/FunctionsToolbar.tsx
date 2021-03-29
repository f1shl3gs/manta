import React, {useMemo, useState} from 'react'
import PromqlToolbarSearch from './PromqlToolbarSearch'
import {
  ComponentSize,
  DapperScrollbars,
  EmptyState,
} from '@influxdata/clockface'

import {functionIdentifierTerms} from 'codemirror-promql/complete/promql.terms'
import FunctionItem from './FunctionItem'

const FunctionsToolbar: React.FC = props => {
  const [search, setSearch] = useState('')
  const filtered = useMemo(() => {
    return functionIdentifierTerms.filter(fn => {
      if (fn.info.indexOf(search) >= 0) {
        return true
      }

      if (fn.label.indexOf(search) >= 0) {
        return true
      }

      return fn.detail.indexOf(search) >= 0
    })
  }, [search])

  console.log('search', search)
  console.log('filterd', filtered.length)

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
        key={fn.label}
        fn={fn}
        testID={fn.label}
        onClickFn={cf => console.log('click', cf.label)}
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
