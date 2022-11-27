import {useCallback, useEffect, useState} from 'react'
import constate from 'constate'
import {useSearchParams} from 'react-router-dom'

const escapeKeyCode = 27
const PRESENTATION_KEY = 'presentation'

const dispatchResizeEvent = () => {
  setTimeout(() => {
    // Uses longer event object creation method due to IE compatibility.
    const ev = document.createEvent('HTMLEvents')
    ev.initEvent('resize', false, true)
    window.dispatchEvent(ev)
  }, 50)
}

const [PresentationModeProvider, usePresentationMode] = constate(() => {
  const [searchParams, setSearchParams] = useSearchParams()
  const [inPresentationMode, setInPresentationMode] = useState(() => {
    const params = new URLSearchParams(window.location.search)
    const defaultValue = params.get(PRESENTATION_KEY)

    return !!(defaultValue && defaultValue === 'true')
  })

  const setPresentation = useCallback(
    (b: boolean) => {
      // History will recorde all urls pushed before, even if the are all the same,
      // in which case when user click `go back`, they will see nothing changed,
      // cause the url is the same
      if (inPresentationMode === b) {
        return
      }

      searchParams.set(PRESENTATION_KEY, b ? 'true' : 'false')
      setSearchParams(searchParams)

      setInPresentationMode(b)
      dispatchResizeEvent()
    },
    [inPresentationMode, searchParams, setSearchParams]
  )

  const toggle = useCallback(() => {
    setPresentation(!inPresentationMode)
  }, [inPresentationMode, setPresentation])

  const escapePresentationMode = useCallback(
    (event: KeyboardEvent) => {
      if (!inPresentationMode) {
        return
      }

      if (event.key === 'Escape' || event.keyCode === escapeKeyCode) {
        setPresentation(false)
      }
    },
    [setPresentation, inPresentationMode]
  )

  useEffect(() => {
    window.addEventListener('keyup', escapePresentationMode)

    // TODO
    // const unListen = history.listen(() => {
    //   setInPresentationMode(false)
    //   dispatchResizeEvent()
    // })

    return () => {
      window.removeEventListener('keyup', escapePresentationMode)
      // unListen()
    }
  })

  return {
    inPresentationMode,
    togglePresentationMode: toggle,
  }
})

export {PresentationModeProvider, usePresentationMode}
