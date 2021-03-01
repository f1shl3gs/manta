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
    <Overlay visible className={'veo-overlay'} onEscape={onDismiss}>
      <div className={'veo'}>
        <CheckProvider id={id}>
          <CheckEditor />
        </CheckProvider>
      </div>
    </Overlay>
  )
}

export default CheckOverlay
