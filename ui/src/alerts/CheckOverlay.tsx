// Libraries
import React from 'react'

// Components
import {Button, Overlay} from '@influxdata/clockface'
import {useHistory, useParams} from 'react-router-dom'
import {CheckProvider} from './useCheck'
import CheckEditor from './CheckEditor'

const CheckOverlay: React.FC = () => {
  const {id} = useParams<{id: string}>()
  const history = useHistory()
  const onDismiss = () => history.goBack()

  return (
    <Overlay visible>
      <Overlay.Container maxWidth={800}>
        <Overlay.Header title={'Check'} onDismiss={onDismiss}>
          <div>hhh</div>
        </Overlay.Header>

        <Overlay.Body>
          <CheckProvider id={id}>
            <CheckEditor />
          </CheckProvider>
        </Overlay.Body>

        <Overlay.Footer>
          <Button text={'Close'} onClick={onDismiss} />
        </Overlay.Footer>
      </Overlay.Container>
    </Overlay>
  )
}

export default CheckOverlay
