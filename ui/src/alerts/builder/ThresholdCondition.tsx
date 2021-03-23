// Libraries
import React from 'react'

// Components
import ThresholdStatement from './ThresholdStatement'
import DashedButton from 'shared/components/DashedButton'
import {ComponentSize} from '@influxdata/clockface'
import ThresholdRangeInput from './ThresholdRangeInput'
import ThresholdValueInput from './ThresholdValueInput'

// Types
import {CheckStatusLevel, Threshold} from '../../types/Check'

// Constants
import {LEVEL_COMPONENT_COLORS} from '../../constants/level'
import {useCheck} from '../checks/useCheck'

interface Props {
  level: CheckStatusLevel
  threshold: Threshold
}

const ThresholdCondition: React.FC<Props> = props => {
  const {level, threshold} = props
  const {onAddCondition} = useCheck()

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
      changeThresholdType={() => console.log('change threshold type')}
    >
      {threshold.type === 'inside' || threshold.type === 'outside' ? (
        <ThresholdRangeInput
          threshold={threshold}
          changeRange={v => console.log('changeRange', v)}
        />
      ) : (
        <ThresholdValueInput
          threshold={threshold}
          changeValue={v => console.log('changeValue', v)}
        />
      )}
    </ThresholdStatement>
  )
}

export default ThresholdCondition
