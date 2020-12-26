import {useCallback, useState} from 'react'

const useVisOptions = () => {
  const [viewing, setViewing] = useState(false)
  const toggleVisOptions = useCallback(() => {
    setViewing(!viewing)
  }, [viewing])

  return {
    toggleVisOptions,
    isViewingVisOptions: viewing,
  }
}
