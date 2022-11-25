// Libraries
import React, {FunctionComponent, useCallback, useState} from 'react'
import {useNavigate} from 'react-router-dom'

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
import useKeyPress from 'src/shared/useKeyPress'
import useFetch from 'src/shared/useFetch'
import {
  defaultErrorNotification,
  defaultSuccessNotification,
  useNotification,
} from 'src/shared/components/notifications/useNotification'
import {useOrganization} from 'src/organizations/useOrganizations'
import {useResources} from 'src/shared/components/GetResources'

// Utils
import {downloadTextFile} from 'src/shared/download'

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
  const [content, setContent] = useState(defaultConfig)
  const navigate = useNavigate()
  const {notify} = useNotification()
  const {id: orgId} = useOrganization()
  const {reload} = useResources()
  const onDismiss = useCallback(() => {
    if (window.history.state.idx > 0) {
      navigate(-1)
    } else {
      const pathname = window.location.pathname.replace('/new', '')
      navigate(pathname)
    }
  }, [navigate])
  const {run: create} = useFetch('/api/v1/configurations', {
    method: 'POST',
    onError: err => {
      notify({
        ...defaultErrorNotification,
        message: `Create new configuration failed, ${err}`,
      })
    },
    onSuccess: _ => {
      notify({
        ...defaultSuccessNotification,
        message: 'Create new configuration success',
      })

      onDismiss()
      reload()
    },
  })

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
            onClick={() => {
              create({
                orgId,
                data: content,
                name: '',
                desc: '',
              })
            }}
          />
        </Overlay.Footer>
      </Overlay.Container>
    </Overlay>
  )
}

export default ConfigurationWizard
