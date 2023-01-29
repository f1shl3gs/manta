// Libraries
import React, {FunctionComponent} from 'react'
import {useParams} from 'react-router-dom'
import {connect, ConnectedProps} from 'react-redux'

// Components
import GetResource from 'src/resources/components/GetResource'
import {Overlay} from '@influxdata/clockface'
import ErrorBoundary from 'src/shared/components/ErrorBoundary'
import NotificationEndpointForm from 'src/notification_endpoints/components/NotificationEndpointForm'
import EndpointOverlayFooter from 'src/notification_endpoints/components/EndpointOverlayFooter'

// Types
import {AppState} from 'src/types/stores'
import {ResourceType} from 'src/types/resources'

// Actions
import {updateNotificationEndpoint} from 'src/notification_endpoints/actions/thunks'

// Selectors
import {getEndpoint} from 'src/notification_endpoints/selectors'

// Hooks
import useEscape from 'src/shared/useEscape'

const mstp = (state: AppState) => {
  return {
    endpoint: getEndpoint(state),
  }
}

const mdtp = {
  onSave: updateNotificationEndpoint,
}

const connector = connect(mstp, mdtp)
type Props = ConnectedProps<typeof connector>

const EditNotificationOverlay: FunctionComponent<Props> = ({onSave}) => {
  const {id} = useParams()

  const onDismiss = useEscape()

  return (
    <Overlay visible={true}>
      <GetResource
        resources={[{id: id, type: ResourceType.NotificationEndpoints}]}
      >
        <Overlay.Container maxWidth={800}>
          <ErrorBoundary>
            <Overlay.Header title={'Edit Notification endpoint'} />

            <Overlay.Body>
              <NotificationEndpointForm />
            </Overlay.Body>

            <EndpointOverlayFooter onSave={onSave} onCancel={onDismiss} />
          </ErrorBoundary>
        </Overlay.Container>
      </GetResource>
    </Overlay>
  )
}

export default connector(EditNotificationOverlay)
