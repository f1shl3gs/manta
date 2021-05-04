// Libraries
import React, {useCallback} from 'react'
import {useHistory} from 'react-router-dom'

// Components
import {Overlay, SpinnerContainer, TechnoSpinner} from '@influxdata/clockface'
import OtclOverlayHeader from './OtclOverlayHeader'
import OtclOverlayContent from './OtclOverlayContent'

// Hooks
import {OtclProvider, useOtcl} from './useOtcl'
import {useOrgID} from '../../shared/useOrg'

const OtclOverlay: React.FC = () => {
  const {loading} = useOtcl()
  const history = useHistory()
  const orgID = useOrgID()
  const onDismiss = useCallback(() => {
    history.push(`/orgs/${orgID}/data/otcls`)
  }, [orgID, history])

  return (
    <Overlay visible className={'veo-overlay'} onEscape={onDismiss}>
      <div className={'veo'}>
        <SpinnerContainer
          loading={loading}
          spinnerComponent={<TechnoSpinner />}
        >
          <OtclOverlayHeader onDismiss={onDismiss} />
          <OtclOverlayContent />
        </SpinnerContainer>
      </div>
    </Overlay>
  )
}

export default () => (
  <OtclProvider>
    <OtclOverlay />
  </OtclProvider>
)
