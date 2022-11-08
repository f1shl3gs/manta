import {Overlay} from '@influxdata/clockface'
import React, {FunctionComponent, useCallback, useEffect} from 'react'
import {useNavigate} from 'react-router-dom'

const ConfigurationWizard: FunctionComponent = () => {
  const navigate = useNavigate()
  const onDismiss = useCallback(() => {
    if (window.history.state.idx > 0) {
      navigate(-1)
    } else {
      const pathname = window.location.pathname.replace('/new', '')
      navigate(pathname)
    }
  }, [navigate])

  // handle esc key press
  useEffect(() => {
    const handleEsc = event => {
      if (event.keyCode === 27) {
        onDismiss()
      }
    }

    window.addEventListener('keydown', handleEsc)

    return () => {
      window.removeEventListener('keydown', handleEsc)
    }
  }, [onDismiss])

  return (
    <Overlay visible={true}>
      <Overlay.Container>
        <Overlay.Header title={'Create Configuration'} onDismiss={onDismiss} />

        <Overlay.Body>blah blah</Overlay.Body>
      </Overlay.Container>
    </Overlay>
  )
}

export default ConfigurationWizard
