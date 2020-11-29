import { RemoteDataState } from '@influxdata/clockface';

const remoteDataState = (loading: boolean, err: any): RemoteDataState => {
  if (err !== undefined) {
    return RemoteDataState.Error;
  }

  if (loading) {
    return RemoteDataState.Loading;
  }

  return RemoteDataState.Done;
};

export default remoteDataState;
