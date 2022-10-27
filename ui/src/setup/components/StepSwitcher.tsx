// Library
import React, {FC} from 'react'

// Components
import {Welcome} from './Welcome'
import {Admin} from './Admin'

// Hooks
import {useStep} from '../useStep'
import {OnboardProvider} from '../useOnboard'

export const StepSwitcher: FC = () => {
  const {step} = useStep()

  switch (step) {
    case 0:
      return <Welcome />
    case 1:
      return (
        <OnboardProvider>
          <Admin />
        </OnboardProvider>
      )
    default:
      return <div />
  }
}
