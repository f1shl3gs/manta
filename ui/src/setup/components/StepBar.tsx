// Libraries
import React, {FC} from 'react'

// Hooks
import {useStep} from '../useStep'
import {ProgressBar} from '@influxdata/clockface'

export const StepBar: FC = () => {
  const {step} = useStep()

  return <ProgressBar max={2} value={step} />
}
