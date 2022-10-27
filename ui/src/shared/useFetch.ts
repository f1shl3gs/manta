import {useCallback, useEffect, useReducer, useRef} from 'react'
import {RemoteDataState} from '@influxdata/clockface'

interface State<T> {
  data?: T
  error?: Error
  loading: RemoteDataState,
  run: () => void
}

type Action<T> =
  | {type: 'loading'}
  | {type: 'fetched'; payload: T}
  | {type: 'error'; payload: Error}

interface RequestOptions<T> {
  method?: string,
  body?: object | string,
  manual?: boolean,

  onSuccess?: (data: T | undefined) => void,
  onError?: (err: Error) => void
}

function useFetch<T = any>(
  url: string,
  options?: RequestOptions<T>,
): State<T> {
  const { method, body, manual, onError, onSuccess} = options ||  {method: 'GET', body: undefined, manual: false}

  // Keep state logic separated
  const fetchReducer = (state: State<T>, action: Action<T>): State<T> => {
    switch (action.type) {
      case 'loading':
        return {...state, loading: RemoteDataState.Loading}
      case 'fetched':
        return {...state, loading: RemoteDataState.Done, data: action.payload}
      case 'error':
        return {...state, loading: RemoteDataState.Error}
      default:
        return state
    }
  }

  const [state, dispatch] = useReducer(fetchReducer, {
    loading: RemoteDataState.NotStarted,
    data: undefined,
    error: undefined,
    // @ts-ignore
    run: () => {console.error('not implement')}
  })

  const run = useCallback(() => {
    dispatch({type: 'loading'})

    const headers = typeof body === 'object' ? {'Content-type': 'application/json'} : undefined

    fetch(url, {
      method,
      headers,
      body: typeof body === 'object' ? JSON.stringify(body) : null,
    })
      .then(resp => {
        const cloned = resp.clone()
        return cloned.json()
          .catch(_ => resp.text())
      })
      .then(data => {
        dispatch({type: 'fetched', payload: data})

        if (onSuccess) {
          onSuccess(data)
        }
      })
      .catch(err => {
        dispatch({type: 'error', payload: err})

        if (onError) {
          onError(err)
        }
      })
  }, [url, body, method, onSuccess, onError])

  useEffect(() => {
    if (manual) {
      return
    }

    run()
  }, [manual, run])

  return {
    ...state,
    run
  }
}

export default useFetch
