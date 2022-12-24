// Libraries
import React, {createRef, FunctionComponent} from 'react'
import classnames from 'classnames'
import {useDispatch, useSelector} from 'react-redux'

// Components
import {ComponentColor, Icon, IconFont, RightClick} from '@influxdata/clockface'
import QueryTabName from 'src/timeMachine/components/queryEditor/QueryTabName'

// Types
import {DashboardQuery} from 'src/types/dashboards'
import {AppState} from 'src/types/stores'

// Actions
import {setActiveQueryIndex, removeQuery} from 'src/timeMachine/actions'

interface Props {
  index: number
  query: DashboardQuery
}

const QueryTab: FunctionComponent<Props> = ({index, query}) => {
  const dispatch = useDispatch()
  const {activeIndex, queries, activeQuery} = useSelector((state: AppState) => {
    const {activeQueryIndex, viewProperties} = state.timeMachine
    const {queries} = viewProperties

    return {
      queries,
      activeIndex: activeQueryIndex,
      activeQuery: queries[activeQueryIndex],
    }
  })

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

          dispatch(setActiveQueryIndex(index))
        }}
        ref={triggerRef}
      >
        {hideButton()}
        <QueryTabName name={activeQuery.name || ''} />
        {queries.length === 1 ? null : (
          <div
            className={'query-tab--close'}
            onClick={() => dispatch(removeQuery(index))}
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
