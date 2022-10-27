// Library
import React, {FC} from 'react'

// Components
import {Button, ComponentColor, ComponentSize} from '@influxdata/clockface'

// Hooks
import {useStep} from '../useStep'

export const Welcome: FC = () => {
  const {next} = useStep()

  return (
    <div className={'wizard--bookend-step'}>
      <div className={'splash-logo primary'} />

      <h3 className={'wizard-step--title'}>Welcome to Manta</h3>

      <h5 className={'wizard-step--sub-title'}>
        Get started in just a few easy steps
      </h5>

      <Button
        color={ComponentColor.Primary}
        text="Get Started"
        size={ComponentSize.Large}
        onClick={next}
      />
    </div>
  )
}
