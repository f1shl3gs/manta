import constate from 'constate'
import {useCallback, useState} from 'react'

const [ViewOptionProvider, useViewOption] = constate(
  () => {
    const [isViewingVisOptions, setIsViewingVisOptions] = useState(false)
    const onToggleVisOptions = useCallback(() => {
      setIsViewingVisOptions(!isViewingVisOptions)
    }, [isViewingVisOptions])

    return {
      isViewingVisOptions,
      onToggleVisOptions,
    }
  },
  value => value
)

export {ViewOptionProvider, useViewOption}
