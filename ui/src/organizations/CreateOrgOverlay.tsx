import {
  Button,
  ButtonType,
  ComponentColor,
  ComponentStatus,
  Form,
  Input,
  Overlay,
} from '@influxdata/clockface'
import React, {FunctionComponent, useCallback, useEffect, useState} from 'react'
import {useNavigate} from 'react-router-dom'

const CreateOrgOverlay: FunctionComponent = () => {
  const [name, setName] = useState('')
  const navigate = useNavigate()
  const onDismiss = useCallback(() => {
    navigate(-1)
  }, [navigate])
  const submitStatus = /^[a-zA-Z0-9]+$/.test(name)
    ? ComponentStatus.Valid
    : ComponentStatus.Disabled
  const orgErrMessage =
    submitStatus === ComponentStatus.Disabled && name !== ''
      ? 'Invalid organization name'
      : ''

  // handle esc
  useEffect(() => {
    const handleEsc = event => {
      if (event.keyCode === 27) {
        onDismiss()
      }
    }

    window.addEventListener('keydown', handleEsc)

    return () => {
      window.removeEventListener('keydown', handleEsc)
    }
  })

  return (
    <Overlay visible={true}>
      <Overlay.Container maxWidth={500}>
        <Overlay.Header
          title={'Create Organization'}
          testID={'create-org-overlay--header'}
          onDismiss={onDismiss}
        />

        <Form onSubmit={() => console.log('create')}>
          <Overlay.Body>
            <Form.Element
              label={'Organization Name'}
              errorMessage={orgErrMessage}
            >
              <Input
                placeholder="Give your organazition a name"
                name="name"
                autoFocus={true}
                value={name}
                onChange={ev => setName(ev.target.value)}
                testID={'create-org-name-input'}
              />
            </Form.Element>
          </Overlay.Body>

          <Overlay.Footer>
            <Button
              text="Cancel"
              color={ComponentColor.Tertiary}
              onClick={onDismiss}
              testID={'create-org-form-cancel'}
            />

            <Button
              text="Create"
              type={ButtonType.Submit}
              color={ComponentColor.Primary}
              testID="create-org-form-create"
              status={submitStatus}
            />
          </Overlay.Footer>
        </Form>
      </Overlay.Container>
    </Overlay>
  )
}

export default CreateOrgOverlay
