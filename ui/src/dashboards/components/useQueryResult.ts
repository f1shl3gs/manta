// Libraries
import {useEffect, useState} from 'react'
import {FromFluxResult, fromRows} from '@influxdata/giraffe'

// Hooks
import {useAutoRefresh} from 'shared/useAutoRefresh'
import {useOrgID} from 'shared/useOrg'
import {useFetch} from 'shared/useFetch'

// Types
import {DashboardQuery} from 'types/Dashboard'

// Utils
import {Row, transformToRows} from 'utils/transform'

const useQueryResult = (
  queries: DashboardQuery[],
  deps?: any[]
): FromFluxResult => {
  // const
  const {get} = useFetch(`/api/v1/query_range`, {})
  const {start, end, step} = useAutoRefresh()
  const orgID = useOrgID()
  const [result, setResult] = useState<FromFluxResult>(() => {
    return {
      table: fromRows([]),
      fluxGroupKeyUnion: [],
    }
  })

  useEffect(() => {
    const set = new Array<Row[]>(queries.length)

    queries.forEach((q, index) => {
      if (q.hidden) {
        return
      }

      get(
        `?query=${encodeURIComponent(
          q.text
        )}&start=${start}&end=${end}&step=${step}&orgID=${orgID}`
      )
        .then(resp => {
          set[index] = transformToRows(resp, q.name || `Query ${index}`)
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
          })
        })
        .catch(err => {
          console.error('query failed', err)
        })
    })
  }, [queries, start, end, step, get, orgID, deps])

  return result
}

export default useQueryResult
