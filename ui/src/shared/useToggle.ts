import {useCallback, useState} from 'react'

function useToggle(initValue?: boolean) {
  const [value, setValue] = useState(initValue ?? false)
  const toggle = useCallback(() => {
    setValue(prev => !prev)
  }, [setValue])

  return [value, toggle] as const
}

export default useToggle
