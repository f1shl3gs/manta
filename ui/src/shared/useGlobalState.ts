import constate from 'constate'
import {useMemo, useState} from 'react'
import useLocalStorage from 'src/shared/useLocalStorage'

const [GlobalStateProvider, usePresentationMode] = constate(
  () => {
    const [presentationMode, setPresentationMode] = useState()
    const [navbarStatus, setNavbarStatus] = useLocalStorage(
      'NavbarStatus',
      false
    )

    return {
      presentationMode,
      setPresentationMode,
      navbarStatus,
      setNavbarStatus,
    }
  },
  value =>
    useMemo(
      () => ({
        presentationMode: value.presentationMode,
        setPresentationMode: value.setPresentationMode,
      }),
      [value.presentationMode, value.setPresentationMode]
    )
)
