import {RemoteDataState} from '@influxdata/clockface'

export interface Configuration {
  id: string
  created: string
  updated: string
  name: string
  desc: string
  data: string

  status: RemoteDataState
}
