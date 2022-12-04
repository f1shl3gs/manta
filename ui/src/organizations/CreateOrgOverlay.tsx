// Libraries
import React, {FunctionComponent, useState} from 'react'
import {connect, ConnectedProps} from 'react-redux'

// Components
import {
  Button,
  ButtonType,
  ComponentColor,
  ComponentStatus,
  Form,
  Input,
  Overlay,
} from '@influxdata/clockface'

// Actions
import {createOrg} from './actions/thunks'

// Hooks
import useEscape from 'src/shared/useEscape'

const mdtp = {
  createOrg,
}

const connector = connect(null, mdtp)

type Props = ConnectedProps<typeof connector>

const CreateOrgOverlay: FunctionComponent<Props> = ({createOrg}) => {
  const [name, setName] = useState('')
  const onDismiss = useEscape()

  const submitStatus = /^[a-zA-Z0-9]+$/.test(name)
    ? ComponentStatus.Valid
    : ComponentStatus.Disabled
  const orgErrMessage =
    submitStatus === ComponentStatus.Disabled && name !== ''
      ? 'Invalid organization name'
      : ''

  return (
    <Overlay visible={true}>
      <Overlay.Container maxWidth={500}>
        <Overlay.Header
          title={'Create Organization'}
          testID={'create-org-overlay--header'}
          onDismiss={onDismiss}
        />

        <Form
          onSubmit={() => {
            createOrg({name, desc: ''})
          }}
        >
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

export default connector(CreateOrgOverlay)
