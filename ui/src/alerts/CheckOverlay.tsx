// Libraries
import React from 'react'

// Components
import {Button, Input, Overlay} from '@influxdata/clockface'
import {useHistory, useParams} from 'react-router-dom'
import {CheckProvider} from './useCheck'

const CheckOverlay: React.FC = () => {
  const history = useHistory()
  const {id} = useParams<{id: string}>()

  return (
    <Overlay visible>
      <Overlay.Container maxWidth={800}>
        <Overlay.Header
          title={'Check'}
          onDismiss={() => {
            history.goBack()
          }}
        >
          <div>hhh</div>
        </Overlay.Header>

        <Overlay.Body>
          <CheckProvider id={id}>
            <div>a</div>
          </CheckProvider>
        </Overlay.Body>

        <Overlay.Footer>
          <Button
            text={'Close'}
            onClick={() => {
              console.log('onClose')
            }}
          />
        </Overlay.Footer>
      </Overlay.Container>
    </Overlay>
  )
}

export default CheckOverlay
