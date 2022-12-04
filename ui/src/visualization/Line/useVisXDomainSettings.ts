// Libraries
import {useCallback, useEffect, useMemo, useState} from 'react'
import {useSelector} from 'react-redux'

// Types
import {AppState} from 'src/types/stores'

export const useVisXDomainSettings = () => {
  const {start, end} = useSelector((state: AppState) => ({
    start: state.autoRefresh.start,
    end: state.autoRefresh.end,
  }))
  const initialXDomain = useMemo(() => [start * 1000, end * 1000], [start, end])
  const [xDomain, setXDomain] = useState(() => [start * 1000, end * 1000])

  useEffect(() => {
    setXDomain([start * 1000, end * 1000])
  }, [start, end])

  const onSetXDomain = useCallback(
    (ns: number[]) => {
      setXDomain(ns)
    },
    [setXDomain]
  )

  const onResetXDomain = useCallback(() => {
    setXDomain(initialXDomain)
  }, [initialXDomain, setXDomain])

  return {
    xDomain,
    onSetXDomain,
    onResetXDomain,
  }
}

// todo: useVisYDomainSettings
