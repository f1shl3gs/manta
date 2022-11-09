import {useEffect} from 'react'

// https://www.w3.org/TR/uievents-key/#named-key-attribute-values
function useKeyPress(key: string, cb: () => void) {
  useEffect(() => {
    const handler = event => {
      if (event.key === key) {
        cb()
      }
    }

    window.addEventListener('keydown', handler)

    return () => {
      window.removeEventListener('keydown', handler)
    }
  }, [key, cb])
}

export default useKeyPress
