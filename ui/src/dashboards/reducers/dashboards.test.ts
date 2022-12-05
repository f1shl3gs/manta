import {DashboardsState} from '../../types/resources'
import {RemoteDataState} from '@influxdata/clockface'
import {DEFAULT_DASHBOARD_SORT_OPTIONS} from '../../constants/dashboard'
import {normalize} from 'normalizr'
import {DashboardEntities} from '../../types/schemas'
import {arrayOfDashboards} from '../../schemas'
import {Dashboard} from 'src/types/dashboards'
import {dashboardsReducer} from './index'
import {setDashboards} from '../actions/creators'

const initialState = (): DashboardsState => ({
  status: RemoteDataState.Done,
  byID: {},
  allIDs: [],
  searchTerm: '',
  current: null,
  sortOptions: DEFAULT_DASHBOARD_SORT_OPTIONS,
})

describe('dashboard reducer', () => {
  it('can set the dashboards', () => {
    const foo = {
      id: 'foo',
      name: 'd1',
      orgID: '1',
      cells: ['foo_1'],
      status: RemoteDataState.Done,
    }

    const schema = normalize<Dashboard, DashboardEntities, string[]>(
      foo,
      arrayOfDashboards
    )

    const got = dashboardsReducer(
      initialState(),
      setDashboards(RemoteDataState.Done, schema)
    )

    console.log('got', got)
  })
})
