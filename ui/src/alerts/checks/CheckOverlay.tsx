// Libraries
import React from 'react'
import {useHistory} from 'react-router-dom'

// Components
import {Overlay, SpinnerContainer, TechnoSpinner} from '@influxdata/clockface'
import CheckEditor from './CheckEditor'
import CheckOverlayHeader from './CheckOverlayHeader'

// Hooks
import {CheckProvider, useCheck} from './useCheck'
import withProvider from '../../utils/withProvider'

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

export default withProvider(CheckProvider, CheckOverlay)
