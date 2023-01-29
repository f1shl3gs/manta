// Libraries
import React, {FunctionComponent, useEffect} from 'react'
import {useDispatch} from 'react-redux'

// Components
import {Overlay} from '@influxdata/clockface'
import ErrorBoundary from 'src/shared/components/ErrorBoundary'
import NotificationEndpointForm from 'src/notification_endpoints/components/NotificationEndpointForm'
import EndpointOverlayFooter from 'src/notification_endpoints/components/EndpointOverlayFooter'

// Hooks
import useEscape from 'src/shared/useEscape'

// Actions
import {createNotificationEndpoint} from 'src/notification_endpoints/actions/thunks'
import {resetCurrentNotificationEndpoint} from '../actions/creators'

const CreateOverlay: FunctionComponent = () => {
  const dispatch = useDispatch()
  const onDismiss = useEscape()
  const onSave = () => {
    dispatch(createNotificationEndpoint())
  }

  useEffect(() => {
    dispatch(resetCurrentNotificationEndpoint())
  }, [dispatch])

  return (
    <Overlay visible={true}>
      <Overlay.Container maxWidth={800}>
        <ErrorBoundary>
          <Overlay.Header
            title={'Create a Notification endpoint'}
            onDismiss={onDismiss}
          />

          <Overlay.Body>
            <NotificationEndpointForm />
          </Overlay.Body>

          <EndpointOverlayFooter onSave={onSave} onCancel={onDismiss} />
        </ErrorBoundary>
      </Overlay.Container>
    </Overlay>
  )
}

export default CreateOverlay
