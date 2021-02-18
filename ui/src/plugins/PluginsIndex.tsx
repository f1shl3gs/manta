import React from 'react'
import {ComponentSize, Page, SquareGrid} from '@influxdata/clockface'

// Constants
import {OTCL_PLUGINS} from 'plugins/constants/Plugins'
import PluginCard from './PluginCard'

const Title = 'OTCL plugins'

const OTCLPluginsIndex: React.FC = () => {
  return (
    <Page titleTag={Title}>
      <Page.Header fullWidth={false}>
        <Page.Title title={'title'} />
      </Page.Header>

      <Page.Contents fullWidth={false} scrollable={true}>
        <SquareGrid cardSize={'170px'} gutter={ComponentSize.Small}>
          {OTCL_PLUGINS.map(item => (
            <PluginCard
              key={item.id}
              id={item.id}
              name={item.name}
              url={item.url}
              image={item.image}
            />
          ))}
        </SquareGrid>
      </Page.Contents>
    </Page>
  )
}

export default OTCLPluginsIndex
