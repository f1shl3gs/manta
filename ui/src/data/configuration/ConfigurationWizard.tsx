import {
  AlignItems,
  Button,
  ComponentColor,
  ComponentSize,
  FlexDirection,
  Heading,
  HeadingElement,
  IconFont,
  JustifyContent,
  Overlay,
  Panel,
} from '@influxdata/clockface'
import React, {FunctionComponent, useCallback, useState} from 'react'
import {useNavigate} from 'react-router-dom'
import YamlMonacoEditor from 'src/shared/components/YamlMonacoEditor'
import useKeyPress from 'src/shared/useKeyPress'

const ConfigurationWizard: FunctionComponent = () => {
  const [content, setContent] = useState('')
  const navigate = useNavigate()
  const onDismiss = useCallback(() => {
    if (window.history.state.idx > 0) {
      navigate(-1)
    } else {
      const pathname = window.location.pathname.replace('/new', '')
      navigate(pathname)
    }
  }, [navigate])

  // handle esc key press
  useKeyPress('Escape', onDismiss)

  return (
    <Overlay visible={true}>
      <Overlay.Container maxWidth={1200}>
        <Overlay.Header
          title={'Create a Configuration'}
          onDismiss={onDismiss}
        />

        <Overlay.Body>
          <Panel style={{marginBottom: '8px'}}>
            <Panel.Body
              direction={FlexDirection.Row}
              alignItems={AlignItems.Center}
              justifyContent={JustifyContent.SpaceBetween}
              size={ComponentSize.ExtraSmall}
            >
              <Heading element={HeadingElement.H4}>name</Heading>

              <Button
                icon={IconFont.Download_New}
                color={ComponentColor.Secondary}
                text={'Download Config'}
                onClick={() => console.log('download config')}
              />
            </Panel.Body>
          </Panel>

          <div className={'config-overlay'}>
            <YamlMonacoEditor content={content} onChange={setContent} />
          </div>
        </Overlay.Body>

        <Overlay.Footer>
          <Button
            color={ComponentColor.Tertiary}
            text={'Cancel'}
            onClick={() => console.log('cancel')}
          />

          <Button color={ComponentColor.Success} text="Save" />
        </Overlay.Footer>
      </Overlay.Container>
    </Overlay>
  )
}

export default ConfigurationWizard
