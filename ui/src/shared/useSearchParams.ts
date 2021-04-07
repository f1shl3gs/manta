import {useCallback, useState} from 'react'
import {useHistory} from 'react-router-dom'
import constate from 'constate'

const [SearchParamsProvider, useSearchParams] = constate(
  () => {
    const {pathname, search} = window.location
    const history = useHistory()
    const [params] = useState<URLSearchParams>(
      () => new URLSearchParams(search)
    )

    const setParams = useCallback(
      (u: (prev: URLSearchParams) => URLSearchParams) => {
        const next = u(params)
        history.push(`${pathname}?${next.toString()}`)
      },
      [history, params, pathname]
    )

    return {
      params,
      setParams,
    }
  },
  value => value
)

export default useSearchParams

export {SearchParamsProvider}
