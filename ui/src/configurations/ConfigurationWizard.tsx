// Libraries
import React, {FunctionComponent, useCallback, useState} from 'react'

// Components
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
import YamlMonacoEditor from 'src/shared/components/editor/YamlMonacoEditor'

// Hooks
import useEscape from 'src/shared/useEscape'
import {useDispatch} from 'react-redux'

// Utils
import {downloadTextFile} from 'src/utils/download'

// Actions
import {createConfig} from 'src/configurations/actions/thunk'

const defaultConfig = `sources:
  # expose process metrics, e.g. cpu, memory and open files
  selfstat:
    type: selfstat
  node:
    type: node_metrics

transforms:
  metrics:
    type: add_tags
    input:
      - selfstat
      - node
    tags:
      host: \${HOSTNAME} # from env

sinks:
  prom:
    type: prometheus_exporter
    inputs:
      - metrics
    endpoint: 0.0.0.0:9100
`

const ConfigurationWizard: FunctionComponent = () => {
  const dispatch = useDispatch()
  const [content, setContent] = useState(defaultConfig)

  // handle esc key press
  const onDismiss = useEscape()

  const handleSave = useCallback(() => {
    dispatch(createConfig('', '', content))
  }, [dispatch, content])

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
                onClick={() => downloadTextFile(content, 'vertex', '.conf')}
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
            onClick={onDismiss}
          />

          <Button
            color={ComponentColor.Success}
            text="Save"
            testID={'create-configuration--button'}
            onClick={handleSave}
          />
        </Overlay.Footer>
      </Overlay.Container>
    </Overlay>
  )
}

export default ConfigurationWizard
