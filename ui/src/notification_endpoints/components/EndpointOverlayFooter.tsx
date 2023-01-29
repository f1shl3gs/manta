// Libraries
import React, {FunctionComponent} from 'react'

// Component
import {
  Button,
  ComponentColor,
  ComponentStatus,
  Overlay,
} from '@influxdata/clockface'

interface Props {
  onSave: () => void
  onCancel: () => void
}

const EndpointOverlayFooter: FunctionComponent<Props> = ({
  onCancel,
  onSave,
}) => {
  const buttonStatus = ComponentStatus.Default

  return (
    <Overlay.Footer>
      <Button
        text={'Cancel'}
        testID={'endpoint-overlay-cancel--button'}
        color={ComponentColor.Tertiary}
        onClick={onCancel}
      />

      <Button
        text={'Save'}
        testID={'endpoint-overlay-save--button'}
        status={buttonStatus}
        color={ComponentColor.Primary}
        onClick={onSave}
      />
    </Overlay.Footer>
  )
}

export default EndpointOverlayFooter
