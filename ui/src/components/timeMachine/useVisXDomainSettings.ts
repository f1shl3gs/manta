import {useAutoRefresh} from 'shared/useAutoRefresh'
import {useCallback, useMemo, useState} from 'react'

export const useVisXDomainSettings = () => {
  const {start, end} = useAutoRefresh()
  const initialXDomain = useMemo(() => [start * 1000, end * 1000], [start, end])
  const [xDomain, setXDomain] = useState(() => [start * 1000, end * 1000])

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
