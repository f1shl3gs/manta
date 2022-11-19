import {DashboardQuery} from 'src/types/Dashboard'
import {useEffect, useState} from 'react'
import {FromFluxResult, fromRows} from '@influxdata/giraffe'
import {Row, transformToRows} from 'src/dashboards/transform'
import {useAutoRefresh} from 'src/shared/useAutoRefresh'

const useQueryResult = (queries: DashboardQuery[]) => {
  const {start, end, step} = useAutoRefresh()
  const [result, setResult] = useState<FromFluxResult>({
    table: fromRows([]),
    fluxGroupKeyUnion: [],
    resultColumnNames: [],
  })
  const [error, setError] = useState('')

  useEffect(() => {
    const set = new Array<Row[]>(queries.length)

    queries.forEach((q, index) => {
      if (q.hidden) {
        return
      }

      if (q.text === '') {
        return
      }

      // Our useFetch implement is not fit here
      fetch(
        `/api/v1/query_range?query=${encodeURIComponent(
          q.text
        )}&start=${start}&end=${end}&step=${step}`
      )
        .then(resp => {
          if (resp.status !== 200) {
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
        })
        .catch(err => {
          setError(err)
        })
    })
  }, [queries, start, end, step])

  return {
    result,
    error,
  }
}

export default useQueryResult
