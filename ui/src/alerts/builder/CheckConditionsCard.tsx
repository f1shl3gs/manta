import React from 'react'
import BuilderCard from '../builderCard/BuilderCard'
import BuilderCardHeader from '../builderCard/BuilderCardHeader'
import BuilderCardBody from '../builderCard/BuilderCardBody'
import {
  AlignItems,
  ComponentSize,
  FlexBox,
  FlexDirection,
} from '@influxdata/clockface'
import ThresholdConditions from './ThresholdConditions'
import {CheckStatusLevel, Condition, Threshold} from '../../types/Check'
import {useCheck} from '../checks/useCheck'

const CheckConditionsCard: React.FC = () => {
  const {conditions} = useCheck()

  return (
    <BuilderCard
      testID={'builder-conditions'}
      className={'alert-builder--card alert-builder--conditions-card'}
    >
      <BuilderCardHeader title={'Conditions'} />

      <BuilderCardBody addPadding={true} autoHideScrollbars={true}>
        <FlexBox
          direction={FlexDirection.Column}
          alignItems={AlignItems.Stretch}
          margin={ComponentSize.Medium}
        >
          <ThresholdConditions
            conditions={conditions.reduce((map, condition) => {
              map[condition.status] = condition
              return map
            }, {} as {[key: string]: Condition})}
          />
        </FlexBox>
      </BuilderCardBody>
    </BuilderCard>
  )
}

export default CheckConditionsCard
