// Libraries
import React, {createRef, FunctionComponent} from 'react'
import classnames from 'classnames'

// Components
import {ComponentColor, Icon, IconFont, RightClick} from '@influxdata/clockface'
import QueryTabName from 'src/timeMachine/QueryTabName'

// Types
import {DashboardQuery} from 'src/types/dashboard'

// Hooks
import {useQueries} from 'src/timeMachine/useTimeMachine'

interface Props {
  index: number
  query: DashboardQuery
}

const QueryTab: FunctionComponent<Props> = ({index, query}) => {
  const {activeIndex, removeQuery, setActiveIndex, queries} = useQueries()
  const activeQuery = queries[activeIndex]

  const triggerRef = createRef<HTMLDivElement>()
  const queryTabClass = classnames('query-tab', {
    'query-tab__active': index === activeIndex,
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

  return (
    <>
      <div
        className={queryTabClass}
        onClick={() => {
          if (index === activeIndex) {
            return
          }

          setActiveIndex(index)
        }}
        ref={triggerRef}
      >
        {hideButton()}
        <QueryTabName name={activeQuery.name || ''} />
        {queries.length === 1 ? null : (
          <div
            className={'query-tab--close'}
            onClick={() => removeQuery(index)}
          >
            <Icon glyph={IconFont.Remove_New} />
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
