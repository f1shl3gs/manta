// Libraries
import React from 'react'

// Components
import ThresholdCondition from './ThresholdCondition'

// Types
import {CheckStatusLevel, Condition} from 'types/Check'

interface Props {
  conditions: {[key: string]: Condition}
}

const CHECK_LEVELS: CheckStatusLevel[] = ['CRIT', 'WARN', 'INFO']

const ThresholdConditions: React.FC<Props> = props => {
  const {conditions} = props

  return (
    <>
      {CHECK_LEVELS.map(level => {
        return (
          <ThresholdCondition
            key={level}
            level={level}
            // @ts-ignore
            threshold={
              conditions[level] ? conditions[level].threshold : undefined
            }
          />
        )
      })}
    </>
  )
}

export default ThresholdConditions
