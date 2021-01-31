// Libraries
import React from 'react'

// Components
import {Page, SpinnerContainer, TechnoSpinner} from '@influxdata/clockface'
import ReactMarkdown from 'react-markdown'
import OTCLPluginsExplainer from './OTCLPluginsExplainer'

// Hooks
import {useParams} from 'react-router-dom'
import {useFetch} from 'use-http'

// Constants
import {OTCL_PLUGINS} from './constants/Plugins'

// Graphics
import placeholderLogo from 'plugins/graphics/placeholderLogo.svg'

// Utils
import remoteDataState from '../utils/rds'

const PluginDetailsView: React.FC = () => {
  const {id} = useParams<{id: string}>()
  const p = OTCL_PLUGINS.find((item) => item.id === id)

  const {name, markdown, image} = p!

  const {data, error, loading} = useFetch(markdown, {}, [])

  return (
    <Page titleTag={`${name}`}>
      <Page.Header fullWidth={false}>
        <Page.Title title={name} />
      </Page.Header>

      <Page.Contents fullWidth={false} scrollable={true}>
        <SpinnerContainer
          loading={remoteDataState(data, error, loading)}
          spinnerComponent={<TechnoSpinner />}
        >
          <div className={'write-data--details'}>
            <div className={'write-data--details-thumbnail'}>
              <img src={image || placeholderLogo} alt={name} />
            </div>

            <div className={'write-data--details-content markdown-format'}>
              <OTCLPluginsExplainer />
              <ReactMarkdown source={data} />
            </div>
          </div>
        </SpinnerContainer>
      </Page.Contents>
    </Page>
  )
}

export default PluginDetailsView
