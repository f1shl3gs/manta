import {RemoteDataState} from '@influxdata/clockface'

export interface Config {
  id: string
  created: string
  updated: string
  name: string
  desc: string
  data: string

  status: RemoteDataState
}
