// Libraries
import React, {useEffect, useState} from 'react'
import {fromRows} from '@influxdata/giraffe'
import classnames from 'classnames'

// Components
import ErrorBoundary from 'shared/components/ErrorBoundary'
import EmptyQueryView from 'shared/components/EmptyQueryView'
import ViewSwitcher from 'shared/components/ViewSwitcher'

// Hooks
import {useViewProperties} from 'shared/useViewProperties'
import {useAutoRefresh} from 'shared/useAutoRefresh'
import {useQueries} from './useQueries'
import {useFetch} from 'use-http'
import {useOrgID} from 'shared/useOrg'

// Utils
import {Row, transformToRows} from 'utils/transform'
import {RemoteDataState} from '@influxdata/clockface'

const TimeMachineVis: React.FC = () => {
  const {viewProperties} = useViewProperties()
  const timeMachineViewClassName = classnames('time-machine--view', {
    'time-machine--view__empty': false,
  })
  const {start, end, step} = useAutoRefresh()
  const {queries} = useQueries()
  const orgID = useOrgID()

  const [queryResults, setQueryResults] = useState(
    () => new Array<Row[]>(queries.length)
  )

  useEffect(() => {
    setQueryResults(new Array<Row[]>(queries.length))
  }, [queries])

  const {get} = useFetch(`/api/v1`)
  useEffect(() => {
    queries.forEach((q, index) => {
      if (q.hidden) {
        return
      }

      get(
        `/query_range?query=${encodeURI(
          q.text
        )}&start=${start}&end=${end}&step=${step}&orgID=${orgID}`
      )
        .then((resp) => {
          // merge
          setQueryResults((prev) => {
            const next = prev
            next[index] = transformToRows(resp, q.name || `Query ${index}`)
            return next
          })
        })
        .catch((err) => {
          console.error('err', err)
        })
    })
  }, [queries, get, orgID])

  const transformer = (results: Row[][]) => {
    const table = fromRows(
      results.flat().sort((a, b) => Number(a['time']) - Number(b['time']))
    )
    const groupKeys = table.columnKeys.filter(
      (key) => key !== 'time' && key !== 'value'
    )

    return {
      table,
      fluxGroupKeyUnion: groupKeys,
    }
  }

  const gr = transformer(queryResults)

  console.log('qrs', queryResults)
  console.log('gr', gr)

  return (
    <div className={timeMachineViewClassName}>
      <ErrorBoundary>
        {/*<ViewLoadingSpinner loading={rds} />*/}
        <EmptyQueryView
          loading={RemoteDataState.Done}
          queries={queries}
          hasResults={gr?.table.length !== 0}
          errorMessage={undefined}
        >
          <ViewSwitcher giraffeResult={gr!} properties={viewProperties} />
        </EmptyQueryView>
      </ErrorBoundary>
    </div>
  )
}

export default TimeMachineVis
