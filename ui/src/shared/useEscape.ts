import {back} from '@lagunovsky/redux-react-router'
import {useCallback, useEffect} from 'react'
import {useDispatch} from 'react-redux'

function useEscape() {
  // useNavigate will cause re-render
  //
  // https://github.com/remix-run/react-router/issues/7634
  const dispatch = useDispatch()
  const goBack = useCallback(() => {
    dispatch(back())
  }, [dispatch])

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
