import {
  Button,
  ButtonType,
  ComponentColor,
  ComponentStatus,
  Form,
  Input,
  Overlay,
} from '@influxdata/clockface'
import React, {FunctionComponent, useCallback, useState} from 'react'
import {useNavigate} from 'react-router-dom'
import useFetch from 'src/shared/useFetch'
import useKeyPress from 'src/shared/useKeyPress'
import {Organization} from 'src/types/organization'
import {
  defaultErrorNotification,
  useNotify,
} from 'src/shared/components/notifications/useNotification'
import {useOrganizations} from 'src/organizations/useOrganizations'

const CreateOrgOverlay: FunctionComponent = () => {
  const [name, setName] = useState('')
  const {refetch} = useOrganizations()
  const navigate = useNavigate()
  const notify = useNotify()
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

  useKeyPress('Escape', onDismiss)

  const {run: createOrg} = useFetch<Organization>(`/api/v1/organizations`, {
    method: 'POST',
    onSuccess: org => {
      refetch()
      navigate(`/orgs/${org?.id}`)
    },
    onError: err => {
      notify({
        ...defaultErrorNotification,
        message: `Create new organization failed, ${err}`,
      })
    },
  })

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
            createOrg({
              name,
            })
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

export default CreateOrgOverlay
