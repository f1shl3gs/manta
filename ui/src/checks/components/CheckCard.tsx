// Libraries
import React, {FunctionComponent, useCallback} from 'react'
import {connect, ConnectedProps} from 'react-redux'
import {useNavigate} from 'react-router-dom'

// Components
import {FlexBox, ResourceCard} from '@influxdata/clockface'
import LastRunStatus from 'src/checks/components/LastRunStatus'
import CheckCardContext from 'src/checks/components/CheckCardContext'

// Types
import {Check} from 'src/types/checks'

// Actions
import {deleteCheck, patchCheck} from 'src/checks/actions/thunks'

// Selectors
import {useOrg} from 'src/organizations/selectors'
import {fromNow} from 'src/shared/utils/duration'

interface OwnProps {
  check: Check
}

const mdtp = {
  patchCheck,
  deleteCheck,
}

const connector = connect(null, mdtp)
type Props = OwnProps & ConnectedProps<typeof connector>

const CheckCard: FunctionComponent<Props> = ({
  check,
  patchCheck,
  deleteCheck,
}) => {
  const navigate = useNavigate()
  const {id: orgID} = useOrg()

  const handleClick = () => {
    navigate(`/orgs/${orgID}/checks/${check.id}/edit`)
  }

  const handleActiveToggle = useCallback(() => {
    const activeStatus = check.activeStatus === 'active' ? 'inactive' : 'active'

    patchCheck(check.id, {
      activeStatus,
    })
  }, [check.id, check.activeStatus, patchCheck])

  const handleDelete = () => {
    deleteCheck(check.id)
  }

  return (
    <ResourceCard
      key={check.id}
      testID={'check-card'}
      contextMenu={
        <CheckCardContext
          activeStatus={check.activeStatus}
          onActiveToggle={handleActiveToggle}
          onDelete={handleDelete}
        />
      }
    >
      <FlexBox>
        <ResourceCard.EditableName
          testID={'check-editable-name'}
          buttonTestID={'check-editable-name--button'}
          inputTestID={'check-editable-name--input'}
          name={check.name}
          onClick={handleClick}
          onUpdate={name => {
            patchCheck(check.id, {name})
          }}
        />

        <LastRunStatus
          lastRunError={check.lastRunError}
          lastRunStatus={check.lastRunStatus}
        />
      </FlexBox>

      <ResourceCard.EditableDescription
        testID={'check-editable-desc'}
        description={check.desc}
        onUpdate={desc => {
          patchCheck(check.id, {desc})
        }}
      />

      <ResourceCard.Meta>
        {`Last Modified: ${fromNow(check.updated)}`}
        {`Last completed: ${check.latestCompleted || 'Never'}`}
      </ResourceCard.Meta>
    </ResourceCard>
  )
}

export default connector(CheckCard)
