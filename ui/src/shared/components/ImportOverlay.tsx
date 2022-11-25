// Libraries
import React, {
  ChangeEvent,
  FunctionComponent,
  useCallback,
  useState,
} from 'react'
import {useNavigate} from 'react-router-dom'

// Components
import {
  Button,
  ButtonType,
  ComponentColor,
  Form,
  Overlay,
  TextArea,
} from '@influxdata/clockface'
import useKeyPress from 'src/shared/useKeyPress'

interface Props {
  resourceName: string
  onSubmit: (string) => void
}

const ImportOverlay: FunctionComponent<Props> = ({resourceName, onSubmit}) => {
  const navigate = useNavigate()
  const [content, setContent] = useState('')

  useKeyPress('Escape', () => {
    navigate(-1)
  })

  const onDismiss = () => {
    navigate(-1)
  }

  const onChange = (ev: ChangeEvent<HTMLTextAreaElement>): void => {
    setContent(ev.target.value)
  }

  const handleSubmit = useCallback(() => {
    onSubmit(content)
  }, [content, onSubmit])

  return (
    <Overlay visible={true} testID={`${resourceName}-import--overlay`}>
      <Overlay.Container maxWidth={800}>
        <Form onSubmit={handleSubmit}>
          <Overlay.Header
            title={`Import ${resourceName}`}
            onDismiss={onDismiss}
          />

          <Overlay.Body>
            <TextArea
              value={content}
              onChange={onChange}
              testID={'import-overlay--textarea'}
            />
          </Overlay.Body>

          <Overlay.Footer>
            <Button
              text={`Import JSON as ${resourceName}`}
              testID={`submit-${resourceName}-button`}
              color={ComponentColor.Primary}
              type={ButtonType.Submit}
            />
          </Overlay.Footer>
        </Form>
      </Overlay.Container>
    </Overlay>
  )
}

export default ImportOverlay
