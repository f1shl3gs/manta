import {useCallback, useEffect, useState} from 'react'
import {useHistory, useLocation} from 'react-router-dom'
import constate from 'constate'

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

const [PresentationModeProvider, usePresentationMode] = constate(
  () => {
    const history = useHistory()
    const location = useLocation()
    const [inPresentationMode, setInPresentationMode] = useState(() => {
      const params = new URLSearchParams(window.location.search)
      const defaultValue = params.get(PRESENTATION_KEY)

      return !!(defaultValue && defaultValue === 'true')
    })

    const setPresentation = useCallback(
      (b: boolean) => {
        const params = new URLSearchParams(location.search)
        params.set(PRESENTATION_KEY, b ? 'true' : 'false')
        history.push(`${location.pathname}?${params.toString()}`)

        setInPresentationMode(b)
        dispatchResizeEvent()
      },
      [history, location]
    )

    const toggle = useCallback(() => {
      console.log('toggle', !inPresentationMode)
      setPresentation(!inPresentationMode)
    }, [inPresentationMode, setPresentation])

    const escapePresentationMode = useCallback(
      event => {
        if (event.key === 'Escape' || event.keyCode === escapeKeyCode) {
          setPresentation(false)
        }
      },
      [setPresentation]
    )

    useEffect(() => {
      window.addEventListener('keyup', escapePresentationMode)
      const unListen = history.listen(() => {
        setInPresentationMode(false)
        dispatchResizeEvent()
      })

      return () => {
        window.removeEventListener('keyup', escapePresentationMode)
        unListen()
      }
    })

    return {
      inPresentationMode,
      togglePresentationMode: toggle,
    }
  },
  value => value
)

export {PresentationModeProvider, usePresentationMode}
