import constate from 'constate'
import {useCallback, useState} from 'react'

const [ManualRefreshProvider, useManualRefresh] = constate(() => {
  const [count, setCount] = useState(0)
  const refresh = useCallback(() => {
    setCount((prevState) => prevState + 1)
  }, [])

  return {
    count,
    refresh,
  }
})

export {ManualRefreshProvider, useManualRefresh}
