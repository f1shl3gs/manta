// Libraries
import React, {FunctionComponent, useEffect, useRef, useState} from 'react'
import {RemoteDataState} from '@influxdata/clockface'
import {useParams} from 'react-router-dom'
import {useSelector} from 'react-redux'
import {fromRows} from '@influxdata/giraffe'

// Components
import EmptyQueryView from 'src/cells/components/EmptyQueryView'
import View from 'src/visualization/View'

// Types
import {ViewProperties} from 'src/types/cells'
import {AppState} from 'src/types/stores'
import {FromFluxResult} from '@influxdata/giraffe'

// Utils
import {executeQuery} from 'src/timeMachine/actions/thunks'

// Hooks
import useIntersectionObserver from 'src/shared/userIntersectionObserver'

interface Props {
  viewProperties: ViewProperties
}

function shallowEqual(object1, object2) {
  const keys1 = Object.keys(object1)
  const keys2 = Object.keys(object2)
  if (keys1.length !== keys2.length) {
    return false
  }

  for (const key of keys1) {
    if (object1[key] !== object2[key]) {
      return false
    }
  }

  return true
}

const TimeSeries: FunctionComponent<Props> = ({viewProperties}) => {
  const {orgID} = useParams()
  const [loading, setLoading] = useState(RemoteDataState.NotStarted)
  const [error, setError] = useState<string>(undefined)
  const [result, setResult] = useState<FromFluxResult>(undefined)
  const {start, end, step} = useSelector((state: AppState) => {
    return state.autoRefresh
  })

  const ref = useRef<HTMLDivElement>(null)
  const entry = useIntersectionObserver(ref, {})

  const [lastSES, setLastSES] = useState({start, end, step})
  const [refreshed, setRefreshed] = useState(true)
  useEffect(() => {
    if (shallowEqual(lastSES, {start, end, step})) {
      setRefreshed(false)
    } else {
      setRefreshed(true)
      setLastSES({start, end, step})
    }
  }, [start, end, step, lastSES])

  const visible = entry?.isIntersecting ?? false

  useEffect(() => {
    if (!visible && !refreshed) {
      return
    }

    if (!refreshed) {
      return
    }

    if (viewProperties.queries.length === 0) {
      return
    }

    setLoading(RemoteDataState.Loading)
    const promises = []
    viewProperties.queries.forEach(query => {
      if (query.hidden) {
        return
      }

      if (query.text.trim() === '') {
        return
      }

      promises.push(
        executeQuery(viewProperties.type, query, orgID, start, end, step)
      )
    })

    Promise.all(promises)
      .then(rows => {
        const table = fromRows(
          rows.flat().sort((a, b) => Number(a['_time']) - Number(b['_time']))
        )

        setResult({
          table,
          fluxGroupKeyUnion: table.columnKeys.filter(
            key => key !== '_time' && key !== '_value'
          ),
          resultColumnNames: [],
        })
        setLoading(RemoteDataState.Done)
      })
      .catch(err => {
        console.error(err)
        setError(err.toString())
        setLoading(RemoteDataState.Error)
      })
  }, [
    refreshed,
    visible,
    setLoading,
    setError,
    start,
    end,
    step,
    orgID,
    viewProperties.queries,
    viewProperties.type,
  ])

  return (
    <div ref={ref} className={'time-series-container'}>
      <EmptyQueryView
        queries={viewProperties.queries}
        hasResults={
          loading === RemoteDataState.Done &&
          result !== undefined &&
          result.table.length !== 0
        }
        loading={loading}
        errorMessage={error}
      >
        <View result={result} properties={viewProperties} />
      </EmptyQueryView>
    </div>
  )
}

export default TimeSeries
