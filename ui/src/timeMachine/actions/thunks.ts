import {GetState} from 'src/types/stores'
import {DashboardQuery} from 'src/types/dashboards'
import {calculateRange} from 'src/shared/actions/autoRefresh'
import request from 'src/utils/request'
import {Row, transformToRows} from 'src/shared/useQueryResult'
import {getOrg} from 'src/organizations/selectors'
import {TimeRange} from 'src/types/timeRanges'
import {fromRows} from '@influxdata/giraffe'
import {setQueryResult} from './index'
import {RemoteDataState} from '@influxdata/clockface'

const executeQuery = async (
  query: DashboardQuery,
  timeRange: TimeRange,
  orgID: string
): Promise<Row[]> => {
  const {start, end, step} = calculateRange(timeRange)

  const resp = await request(`/api/v1/query_range`, {
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
      viewProperties: {queries},
    } = state.timeMachine

    const promises = []
    queries.forEach(query => {
      if (query.hidden) {
        return
      }

      if (query.text === '') {
        return
      }

      promises.push(executeQuery(query, timeRange, org.id))
    })

    const rows = await Promise.all(promises)
    const table = fromRows(
      rows.flat().sort((a, b) => Number(a['time']) - Number(b['time']))
    )

    dispatch(
      setQueryResult(RemoteDataState.Done, {
        table,
        fluxGroupKeyUnion: table.columnKeys.filter(
          key => key !== 'time' && key !== 'value'
        ),
        resultColumnNames: [],
      })
    )
  }
