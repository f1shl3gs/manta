// Libraries
import React from 'react'

// Components
import {
  AlignItems,
  Button,
  ComponentColor,
  ComponentSize,
  FlexBox,
  FlexDirection,
  IconFont,
  JustifyContent,
  ResourceCard,
  SlideToggle,
} from '@influxdata/clockface'
import LastRunStatus from './LastRunStatus'

// Utils
import {relativeTimestampFormatter} from '../../utils/relativeTimestampFormatter'

// Types
import {Check} from '../../types/Check'
import {useChecks} from './useChecks'
import {useHistory} from 'react-router-dom'
import {useOrgID} from '../../shared/useOrg'

interface Props {
  check: Check
}

const CheckCard: React.FC<Props> = props => {
  const {
    check: {
      id,
      name,
      desc,
      updated,
      status,
      latestCompleted,
      lastRunStatus,
      lastRunError,
    },
  } = props
  const {del, reload} = useChecks()
  const orgID = useOrgID()
  const history = useHistory()

  const contextMenu = () => (
    <Button
      icon={IconFont.Trash}
      text={'Delete'}
      color={ComponentColor.Danger}
      size={ComponentSize.ExtraSmall}
      onClick={() => {
        del(id)
          .then(() => {
            reload()
          })
          .catch(err => {
            console.error(err)
          })
      }}
    />
  )

  return (
    <ResourceCard
      key={`check-id--${id}`}
      disabled={status === 'inactive'}
      direction={FlexDirection.Row}
      alignItems={AlignItems.Center}
      margin={ComponentSize.Large}
      contextMenu={contextMenu()}
    >
      <FlexBox
        direction={FlexDirection.Column}
        justifyContent={JustifyContent.Center}
        margin={ComponentSize.Medium}
        alignItems={AlignItems.FlexStart}
      >
        <SlideToggle
          active={status !== 'inactive'}
          size={ComponentSize.ExtraSmall}
          onChange={() => console.log('toggle')}
          style={{flexBasis: '16px'}}
        />
        <LastRunStatus
          lastRunStatus={lastRunStatus}
          lastRunError={lastRunError}
        />
      </FlexBox>

      <FlexBox
        direction={FlexDirection.Column}
        margin={ComponentSize.Small}
        alignItems={AlignItems.FlexStart}
      >
        <ResourceCard.EditableName
          name={name}
          noNameString={'Name this Check'}
          onClick={() => {
            history.push(`/orgs/${orgID}/alerts/checks/${id}`)
          }}
          onUpdate={v => console.log('v')}
        />
        <ResourceCard.EditableDescription
          description={desc || ''}
          placeholder={`Describe ${name}`}
          onUpdate={v => console.log('onUpdate')}
        />
        <ResourceCard.Meta>
          <>Last completed at {latestCompleted}</>
          <>{relativeTimestampFormatter(updated, 'Last updated ')}</>
        </ResourceCard.Meta>
      </FlexBox>
    </ResourceCard>
  )
}

export default CheckCard
