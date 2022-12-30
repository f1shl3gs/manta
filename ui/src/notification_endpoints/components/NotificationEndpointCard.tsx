// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {ResourceCard} from '@influxdata/clockface'

// Types
import {NotificationEndpoint} from 'src/types/notificationEndpoints'

// Utils
import {fromNow} from 'src/shared/utils/duration'

interface OwnProps {
  notificationEndpoint: NotificationEndpoint
}

const NotificationEndpointCard: FunctionComponent<OwnProps> = ({
  notificationEndpoint,
}) => {
  return (
    <ResourceCard>
      <ResourceCard.EditableName
        name={notificationEndpoint.name}
        onClick={() => console.log('todo')}
        onUpdate={v => console.log(v)}
      />

      <ResourceCard.EditableDescription
        description={notificationEndpoint.desc}
        onUpdate={d => console.log(d)}
      />

      <ResourceCard.Meta>
        <>Modified: {fromNow(notificationEndpoint.updated)}</>
      </ResourceCard.Meta>
    </ResourceCard>
  )
}

export default NotificationEndpointCard
