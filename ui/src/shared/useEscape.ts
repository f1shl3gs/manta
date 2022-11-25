import {useCallback, useEffect} from 'react'
import {useNavigate} from 'react-router-dom'

function useEscape() {
  const navigate = useNavigate()
  const goBack = useCallback(() => {
    navigate(-1)
  }, [navigate])

  useEffect(() => {
    const handler = event => {
      if (event.key === 'Escape') {
        goBack()
      }
    }

    window.addEventListener('keydown', handler)

    return () => {
      window.removeEventListener('keydown', handler)
    }
  }, [goBack])

  return goBack
}

export default useEscape
