// Libraries
import React from 'react'

// Components
import ThresholdStatement from './ThresholdStatement'
import DashedButton from 'shared/components/DashedButton'
import {ComponentSize} from '@influxdata/clockface'
import ThresholdRangeInput from './ThresholdRangeInput'
import ThresholdValueInput from './ThresholdValueInput'

// Types
import {
  CheckStatusLevel,
  GreatThanThreshold,
  InsideThreshold,
  OutsideThreshold,
  Threshold,
} from 'types/Check'

// Constants
import {LEVEL_COMPONENT_COLORS} from '../../../constants/level'
import {useCheck} from '../useCheck'

interface Props {
  level: CheckStatusLevel
  threshold: Threshold
}

const ThresholdCondition: React.FC<Props> = props => {
  const {level, threshold} = props
  const {onAddCondition, onChangeCondition} = useCheck()

  if (!threshold) {
    return (
      <DashedButton
        text={`+ ${level}`}
        color={LEVEL_COMPONENT_COLORS[level]}
        size={ComponentSize.Large}
        onClick={() => onAddCondition(level)}
        testID={`add-threshold-condition-${level}`}
      />
    )
  }

  return (
    <ThresholdStatement
      level={level}
      threshold={threshold}
      // removeLevel={() => console.log('rl')}
      changeThresholdType={type =>
        onChangeCondition({
          status: level as CheckStatusLevel,
          threshold: {
            type,
            value: 0,
          } as GreatThanThreshold,
        })
      }
    >
      {threshold.type === 'inside' || threshold.type === 'outside' ? (
        <ThresholdRangeInput
          threshold={threshold}
          changeRange={(min, max) => {
            onChangeCondition({
              status: level as CheckStatusLevel,
              threshold: {
                type: threshold.type,
                min,
                max,
              },
            })
          }}
        />
      ) : (
        <ThresholdValueInput
          threshold={threshold}
          changeValue={v =>
            onChangeCondition({
              status: level as CheckStatusLevel,
              threshold: {
                type: threshold.type,
                value: v,
              },
            })
          }
        />
      )}
    </ThresholdStatement>
  )
}

export default ThresholdCondition
