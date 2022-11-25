import {DashboardQuery} from 'src/types/Dashboard'
import {useEffect, useState} from 'react'
import {FromFluxResult, fromRows} from '@influxdata/giraffe'
import {useAutoRefresh} from 'src/shared/useAutoRefresh'
import {useParams} from 'react-router-dom'
import {RemoteDataState} from '@influxdata/clockface'

export type Result = {
  metric: {
    [key: string]: string
  }
  values: [[number, string]]
}

export type PromResp = {
  status: string
  data: {
    resultType: string
    result: Result[]
  }
}

export interface Row {
  [key: string]: string | number
}

export const transformToRows = (resp: PromResp): Row[] => {
  if (!resp) {
    return []
  }

  if (resp.status !== 'success') {
    return []
  }

  return resp.data.result
    .map((item: Result) => {
      const {metric, values} = item

      return values.map(val => {
        return {
          ...metric,
          time: val[0] * 1000,
          value: Number(val[1]),
        }
      })
    })
    .flat()
}

export const transformPromResp = (
  resp?: PromResp
): FromFluxResult | undefined => {
  if (!resp) {
    return undefined
  }

  if (resp.status !== 'success') {
    return undefined
  }

  const rows = resp.data.result
    .map((item: Result) => {
      const {metric, values} = item

      return values.map(val => {
        return {
          ...metric,
          time: val[0] * 1000,
          value: Number(val[1]),
        }
      })
    })
    .flat()

  const table = fromRows(rows)
  const groupKeys = table.columnKeys.filter(
    key => key !== 'time' && key !== 'value'
  )

  return {
    table,
    fluxGroupKeyUnion: groupKeys,
    // TODO:
    resultColumnNames: [],
  }
}

const useQueryResult = (queries: DashboardQuery[]) => {
  const [error, setError] = useState('')
  const {start, end, step} = useAutoRefresh()
  const {orgId} = useParams()
  const [loading, setLoading] = useState(RemoteDataState.NotStarted)
  const [result, setResult] = useState<FromFluxResult>({
    table: fromRows([]),
    fluxGroupKeyUnion: [],
    resultColumnNames: [],
  })

  useEffect(() => {
    const set = new Array<Row[]>(queries.length)

    queries.forEach((q, index) => {
      if (q.hidden) {
        setLoading(RemoteDataState.Done)
        return
      }

      if (q.text === '') {
        setLoading(RemoteDataState.Done)
        return
      }

      setLoading(RemoteDataState.Loading)
      // Our useFetch implement is not fit here
      fetch(
        `/api/v1/query_range?query=${encodeURIComponent(
          q.text
        )}&start=${start}&end=${end}&step=${step}&orgId=${orgId}`
      )
        .then(resp => {
          if (resp.status !== 200) {
            setLoading(RemoteDataState.Error)
            setError(resp.statusText)
            return
          }

          return resp.json()
        })
        .then(resp => {
          set[index] = transformToRows(resp)
          const table = fromRows(
            set.flat().sort((a, b) => {
              return Number(a['time']) - Number(b['time'])
            })
          )

          setResult({
            table,
            fluxGroupKeyUnion: table.columnKeys.filter(
              key => key !== 'time' && key !== 'value'
            ),
            resultColumnNames: [],
          })
          setLoading(RemoteDataState.Done)
        })
        .catch(err => {
          setError(err)
          setLoading(RemoteDataState.Error)
        })
    })
  }, [queries, start, end, step, orgId])

  return {
    result,
    loading,
    error,
  }
}

export default useQueryResult
