import {RemoteDataState} from '@influxdata/clockface'

/*
 * @deprecated
 * */
const remoteDataState = (data: any, error: any, loading: boolean) => {
  if (loading) {
    return RemoteDataState.Loading
  }

  if (error !== undefined) {
    return RemoteDataState.Error
  }

  if (data !== undefined) {
    return RemoteDataState.Done
  } else {
    return RemoteDataState.Loading
  }
}

export default remoteDataState
