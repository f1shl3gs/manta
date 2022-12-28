// Libraries
import React, {FunctionComponent, useState} from 'react'
import {useParams} from 'react-router-dom'
import {useDispatch} from 'react-redux'

// Components
import {
  Alert,
  Button,
  ComponentColor,
  ComponentStatus,
  Form,
  IconFont,
  Input,
  Overlay,
} from '@influxdata/clockface'

// Hooks
import useEscape from 'src/shared/useEscape'

// Actions
import {upsertSecret} from 'src/secrets/actions/thunks'

const warningText =
  'Updating this secret could cause queries that rely on this secret to break'

const EditSecretOverlay: FunctionComponent = () => {
  const {key} = useParams()
  const [value, setValue] = useState('')
  const onDismiss = useEscape()
  const dispatch = useDispatch()

  return (
    <Overlay visible={true}>
      <Overlay.Container maxWidth={800}>
        <Overlay.Header title={'Edit Secret'} />

        <Overlay.Body>
          <Alert
            color={ComponentColor.Warning}
            className={'alert'}
            icon={IconFont.AlertTriangle}
          >
            {warningText}
          </Alert>

          <br />

          <Form.Label label={'Value'} />
          <Input
            autoFocus={true}
            required={true}
            testID={'secret-value--input'}
            name={'value'}
            placeholder={'your_secret_value'}
            titleText={
              'This is the value that will be injected by the server when your secret is in use'
            }
            value={value}
            onChange={ev => setValue(ev.target.value)}
          />
        </Overlay.Body>

        <Overlay.Footer>
          <Button
            text={'Cancel'}
            color={ComponentColor.Tertiary}
            onClick={onDismiss}
          />

          <Button
            text={'Update'}
            testID={'update-secret--button'}
            onClick={() => dispatch(upsertSecret({key, value}))}
            color={ComponentColor.Success}
            status={
              value.trim() !== ''
                ? ComponentStatus.Default
                : ComponentStatus.Disabled
            }
          />
        </Overlay.Footer>
      </Overlay.Container>
    </Overlay>
  )
}

export default EditSecretOverlay
