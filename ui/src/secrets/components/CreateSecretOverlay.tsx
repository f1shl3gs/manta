// Libraries
import React, {FunctionComponent, useState} from 'react'
import {useDispatch, useSelector} from 'react-redux'

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

// Selectors
import {getAllSecrets} from 'src/secrets/selectors'

const warningText =
  'Make sure you know your secret value! You will be able to reference the secret in queries by key but you will not be able to see the value again.'

const CreateSecretOverlay: FunctionComponent = () => {
  const [secret, setSecret] = useState({key: '', value: ''})
  const onDismiss = useEscape()
  const dispatch = useDispatch()
  const secrets = useSelector(getAllSecrets)
  const handleKeyValidation = (key: string): string | null => {
    if (!key) {
      return null
    }

    if (!/^[a-zA-Z0-9]+$/.test(key)) {
      return 'Only [a-zA-Z0-9]+$ allowd'
    }

    if (key.trim() === '') {
      return 'Key is required'
    }
    const existingKeys = secrets.map(s => s.key)

    if (existingKeys.includes(key)) {
      return 'Key is already in use'
    }

    return null
  }

  const handleChangeInput = ({target}) => {
    const {name, value} = target

    setSecret(prevState => ({...prevState, [name]: value}))
  }

  const isFormValid = (): boolean => {
    return (
      handleKeyValidation(secret.key) === null && secret.value.trim().length > 0
    )
  }

  return (
    <Overlay visible={true}>
      <Overlay.Container maxWidth={800}>
        <Overlay.Header title={'Create Secret'} />

        <Overlay.Body>
          <Alert
            className="alert"
            icon={IconFont.AlertTriangle}
            color={ComponentColor.Warning}
          >
            {warningText}
          </Alert>

          <br />

          <Form.ValidationElement
            validationFunc={handleKeyValidation}
            label={'Key'}
            value={secret.key}
          >
            {status => (
              <Input
                autoFocus={true}
                testID={'secret-name--input'}
                name={'key'}
                titleText={'This is how you will reference your secret'}
                value={secret.key}
                status={status}
                onChange={handleChangeInput}
              />
            )}
          </Form.ValidationElement>

          <Form.Label label={'value'} />

          <Input
            name={'value'}
            required={true}
            placeholder={'your_secret_value'}
            titleText={
              'This is the value that will be injected by the server when your secret is in use'
            }
            value={secret.value}
            onChange={handleChangeInput}
            testID={'secret-value--input'}
          />
        </Overlay.Body>

        <Overlay.Footer>
          <Button
            text={'Cancel'}
            color={ComponentColor.Tertiary}
            onClick={onDismiss}
          />

          <Button
            text={'Add secret'}
            testID={'create-secret--button'}
            onClick={() => dispatch(upsertSecret(secret))}
            color={ComponentColor.Success}
            status={
              isFormValid() ? ComponentStatus.Default : ComponentStatus.Disabled
            }
          />
        </Overlay.Footer>
      </Overlay.Container>
    </Overlay>
  )
}

export default CreateSecretOverlay
