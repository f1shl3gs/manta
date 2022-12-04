// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {Button, ComponentColor, Overlay, TextArea} from '@influxdata/clockface'
import CopyToClipboard from 'src/shared/components/CopyToClipboard'

// Hooks
import useEscape from '../useEscape'
import {useDispatch} from 'react-redux'

// Utils
import {downloadTextFile} from 'src/utils/download'

// Actions
import {notify} from 'src/shared/actions/notifications'

// Constants
import {defaultSuccessNotification} from 'src/constants/notification'

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
  const dispatch = useDispatch()

  const onCopy = (_copiedText: string, _copyWasSuccessful: boolean): void => {
    dispatch(
      notify({
        ...defaultSuccessNotification,
        message: `Copy to Clipboard successful`,
      })
    )
  }
  const download = () => {
    downloadTextFile(content, name, '', 'application/json')
  }

  const onDismiss = useEscape()

  return (
    <Overlay visible={true} testID={`${resourceName}-export--overlay`}>
      <Overlay.Container maxWidth={900}>
        <Overlay.Header
          title={`Export ${resourceName}`}
          onDismiss={onDismiss}
        />

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
