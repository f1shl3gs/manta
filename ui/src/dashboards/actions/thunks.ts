import request from '../../utils/request';
import {GetState} from '../../types/stores';
import {Dispatch} from 'react'

import {Action, setDashboards} from 'src/dashboards/actions/creators'
import {RemoteDataState} from '@influxdata/clockface';
import {getStatus} from '../../resources/selectors';
import {ResourceType} from '../../types/resources';
import {getOrg} from '../../organizations/selectors';
import { normalize } from 'normalizr';
import { DashboardEntities } from 'src/types/schemas';
import { Dashboard } from 'src/types/Dashboard';
import { arrayOfDashboards } from 'src/schemas';

export const getDashboards = () =>
async (dispatch: Dispatch<Action>, getState: GetState): Promise<Dashboard[]> => {
  try {
    const state = getState()
    if (getStatus(state, ResourceType.Dashboards) === RemoteDataState.NotStarted) {
      dispatch(setDashboards(RemoteDataState.Loading))
    }

    const org = getOrg(state)

    const resp = await request(`/api/v1/dashboards?orgID=${org.id}`)
    if (resp.status !== 200) {
      throw new Error(resp.data.message)
    }

    const dashbaords = normalize<Dashboard, DashboardEntities, string[]>(
      resp.data,
      arrayOfDashboards
    )

    setDashboards(RemoteDataState.Done, )
  }
}