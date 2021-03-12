// Libraries
import React from 'react'

// Components
import {
  Button,
  Overlay,
  SpinnerContainer,
  TechnoSpinner,
} from '@influxdata/clockface'
import {useHistory, useParams} from 'react-router-dom'
import {CheckProvider, useCheck} from './useCheck'
import CheckEditor from './CheckEditor'
import CheckOverlayHeader from './CheckOverlayHeader'

const CheckOverlay: React.FC = () => {
  const history = useHistory()
  const onDismiss = () => history.goBack()
  const {remoteDataState} = useCheck()

  return (
    <Overlay visible className={'veo-overlay'} onEscape={onDismiss}>
      <div className={'veo'}>
        <SpinnerContainer
          loading={remoteDataState}
          spinnerComponent={<TechnoSpinner />}
        >
          <CheckOverlayHeader />

          <CheckEditor />
        </SpinnerContainer>
      </div>
    </Overlay>
  )
}

export default () => {
  const {id} = useParams<{id: string}>()

  return (
    <CheckProvider id={id}>
      <CheckOverlay />
    </CheckProvider>
  )
}
