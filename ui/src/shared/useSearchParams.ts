import {useHistory, useLocation} from 'react-router-dom'
import {useCallback, useMemo} from 'react'

type setStateAction = (prev: URLSearchParams) => URLSearchParams

function useSearchParams(): {
  params: URLSearchParams
  setParams: (u: setStateAction) => void
} {
  const {search, pathname} = useLocation()
  const history = useHistory()
  const params = useMemo(() => new URLSearchParams(search), [search])

  const setParams = useCallback(
    (u: setStateAction) => {
      const next = u(params)
      history.push(`${pathname}?${next.toString()}`)
    },
    [history, params, pathname]
  )

  return {
    params,
    setParams,
  }
}

export default useSearchParams
