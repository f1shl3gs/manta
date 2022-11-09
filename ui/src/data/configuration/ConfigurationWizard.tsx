import {Overlay} from '@influxdata/clockface'
import React, {FunctionComponent, useCallback} from 'react'
import {useNavigate} from 'react-router-dom'
import useKeyPress from 'src/shared/useKeyPress'

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
  useKeyPress('Escape', onDismiss)

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
