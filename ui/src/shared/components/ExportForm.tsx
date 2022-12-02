// Libraries
import React, {FunctionComponent} from 'react'
import {useNavigate} from 'react-router-dom'
import CopyToClipboard from 'src/shared/components/CopyToClipboard'

// Components
import {Button, ComponentColor, Overlay, TextArea} from '@influxdata/clockface'

// Hooks
import useKeyPress from 'src/shared/useKeyPress'
import {
  defaultSuccessNotification,
  useNotify,
} from 'src/shared/components/notifications/useNotification'

// Utils
import {downloadTextFile} from 'src/utils/download'

interface Props {
  resourceName: string
  name: string
  content: string
}

const ExportOverlay: FunctionComponent<Props> = ({
  resourceName,
  name,
  content,
}) => {
  const navigate = useNavigate()
  const notify = useNotify()

  const onDimiss = () => navigate(-1)
  const onCopy = (_copiedText: string, _copyWasSuccessful: boolean): void => {
    notify({
      ...defaultSuccessNotification,
      message: `Copy to Clipboard successful`,
    })
  }
  const download = () => {
    downloadTextFile(content, name, '', 'application/json')
  }

  useKeyPress('Escape', onDimiss)

  return (
    <Overlay visible={true} testID={`${resourceName}-export--overlay`}>
      <Overlay.Container maxWidth={900}>
        <Overlay.Header title={`Export ${resourceName}`} onDismiss={onDimiss} />

        <Overlay.Body>
          <TextArea
            value={content}
            testID={'export-overlay--textarea'}
            readOnly={true}
          />
        </Overlay.Body>

        <Overlay.Footer>
          <Button
            text={'Download JSON'}
            testID={'export-overlay--download'}
            color={ComponentColor.Primary}
            onClick={download}
          />

          <CopyToClipboard text={content} onCopy={onCopy}>
            <Button
              text={'Copy to clipboard'}
              testID={'export-overlay--download'}
              color={ComponentColor.Primary}
            />
          </CopyToClipboard>
        </Overlay.Footer>
      </Overlay.Container>
    </Overlay>
  )
}

export default ExportOverlay
