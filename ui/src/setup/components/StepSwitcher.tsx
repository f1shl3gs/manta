// Library
import React, {FC} from 'react'

// Components
import {Welcome} from 'src/setup/components/Welcome'
import {Admin} from 'src/setup/components/Admin'

// Hooks
import {useStep} from 'src/setup/useStep'
import {OnboardProvider} from 'src/setup/useOnboard'

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
