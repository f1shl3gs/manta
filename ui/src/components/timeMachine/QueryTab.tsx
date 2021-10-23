import React, {createRef} from 'react'
import classnames from 'classnames'
import {DashboardQuery} from '../../types/Dashboard'
import {ComponentColor, Icon, IconFont, RightClick} from '@influxdata/clockface'
import QueryTabName from './QueryTabName'
import {useQueries} from './useQueries'

interface Props {
  queryIndex: number
  query: DashboardQuery
}

const QueryTab: React.FC<Props> = props => {
  const {query, queryIndex} = props
  const {activeIndex, removeQuery, setActiveIndex, queries} = useQueries()
  const activeQuery = queries[activeIndex]

  const triggerRef = createRef<HTMLDivElement>()
  const queryTabClass = classnames('query-tab', {
    'query-tab__active': queryIndex === activeIndex,
    'query-tab__hidden': query.hidden,
  })
  const hideButton = () => {
    const icon = query.hidden ? IconFont.EyeClosed : IconFont.EyeOpen
    return (
      <div
        className={'query-tab--hide'}
        onClick={() => console.log('toggle view')}
      >
        <Icon glyph={icon} />
      </div>
    )
  }

  const removeButton = () => {
    if (queries.length === 1) {
      return null
    }

    return (
      <div
        className={'query-tab--close'}
        onClick={() => removeQuery(queryIndex)}
      >
        <Icon glyph={IconFont.Remove} />
      </div>
    )
  }

  return (
    <>
      <div
        className={queryTabClass}
        onClick={() => {
          if (queryIndex === activeIndex) {
            return
          }

          setActiveIndex(queryIndex)
        }}
        ref={triggerRef}
      >
        {hideButton()}
        <QueryTabName name={activeQuery.name || ''} />
        {queries.length === 1 ? null : (
          <div
            className={'query-tab--close'}
            onClick={() => removeQuery(queryIndex)}
          >
            <Icon glyph={IconFont.Remove} />
          </div>
        )}
      </div>

      <RightClick triggerRef={triggerRef} color={ComponentColor.Primary}>
        <RightClick.MenuItem
          onClick={() => console.log('handle edit active query name')}
        >
          Edit
        </RightClick.MenuItem>

        <RightClick.MenuItem
          onClick={() => console.log('handle remove')}
          disabled={false}
        >
          Remove
        </RightClick.MenuItem>
      </RightClick>
    </>
  )
}

export default QueryTab
