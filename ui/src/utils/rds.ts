import {RemoteDataState} from '@influxdata/clockface'

const remoteDataState = (data: any, error: any, loading: boolean) => {
  if (loading) {
    return RemoteDataState.Loading;
  }

  if (error !== undefined) {
    return RemoteDataState.Error;
  }

  if (data !== undefined) {
    return RemoteDataState.Done;
  }

  return RemoteDataState.NotStarted;
};

export default remoteDataState
