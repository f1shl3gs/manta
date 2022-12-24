// Libraries
import React, {FunctionComponent} from 'react'
import {connect, ConnectedProps} from 'react-redux'

// Components
import ConditionCard from 'src/checks/components/checkBuilder/ConditionCard'
import BuilderCard from 'src/checks/components/builderCard/BuilderCard'
import BuilderCardHeader from 'src/checks/components/builderCard/BuilderCardHeader'
import BuilderCardBody from 'src/checks/components/builderCard/BuilderCardBody'

// Types
import {AppState} from 'src/types/stores'
import {ConditionStatus} from 'src/types/checks'

const CHECK_CONDITION_STATUSES: ConditionStatus[] = [
  'crit',
  'warn',
  'info',
  'ok',
]

const mstp = (state: AppState) => {
  const {conditions} = state.checkBuilder

  return {
    conditions,
  }
}

const connector = connect(mstp, null)
type Props = ConnectedProps<typeof connector>

const ConditionList: FunctionComponent<Props> = ({conditions}) => {
  return (
    <BuilderCard className="alert-builder--card alert-builder--conditions-card">
      <BuilderCardHeader title={'Conditions'} testID={'builder-card-header'} />
      <BuilderCardBody autoHideScrollbars={true} addPadding={true}>
        {CHECK_CONDITION_STATUSES.map(status => (
          <ConditionCard
            key={status}
            condition={{
              status,
              pending: '0s',
              ...conditions[status],
              threshold: conditions[status]
                ? conditions[status].threshold
                : undefined,
            }}
          />
        ))}
      </BuilderCardBody>
    </BuilderCard>
  )
}

export default connector(ConditionList)
