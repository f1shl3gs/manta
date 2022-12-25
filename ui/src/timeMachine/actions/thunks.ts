// Libraries
import {fromRows} from '@influxdata/giraffe'

// Types
import {GetState} from 'src/types/stores'
import {DashboardQuery} from 'src/types/dashboards'
import {RemoteDataState} from '@influxdata/clockface'
import {ViewType} from 'src/types/cells'

// Actions
import {calculateRange} from 'src/shared/actions/autoRefresh'
import {setQueryResult} from 'src/timeMachine/actions'

// Utils
import request from 'src/shared/utils/request'
import {Row, transformToRows} from 'src/shared/utils/transform'
import {getOrg} from 'src/organizations/selectors'

export const executeQuery = async (
  type: ViewType,
  query: DashboardQuery,
  orgID: string,
  start: number,
  end: number,
  step: number
): Promise<Row[]> => {
  const url =
    type === 'single-stat' || type === 'gauge'
      ? `/api/v1/query`
      : `/api/v1/query_range`
  const resp = await request(url, {
    query: {
      start: `${start}`,
      end: `${end}`,
      step: `${step}`,
      orgID,
      query: query.text,
    },
  })
  if (resp.status !== 200) {
    throw new Error(resp.data.message)
  }

  return transformToRows(resp.data)
}

export const loadView =
  () =>
  async (dispatch, getState: GetState): Promise<void> => {
    const state = getState()
    const org = getOrg(state)
    const {
      timeRange,
      viewProperties: {type, queries},
    } = state.timeMachine
    const {start, end, step} = calculateRange(timeRange)

    const promises = []
    queries.forEach(query => {
      if (query.hidden) {
        return
      }

      if (query.text.trim() === '') {
        return
      }

      promises.push(executeQuery(type, query, org.id, start, end, step))
    })

    const rows = await Promise.all(promises)
    const table = fromRows(
      rows.flat().sort((a, b) => Number(a['_time']) - Number(b['_time']))
    )

    dispatch(
      setQueryResult(RemoteDataState.Done, {
        table,
        fluxGroupKeyUnion: table.columnKeys.filter(
          key => key !== '_time' && key !== '_value'
        ),
        resultColumnNames: [],
      })
    )
  }
