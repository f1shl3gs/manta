// Library
import React, {FC} from 'react'

// Components
import {Button, ComponentColor, ComponentSize} from '@influxdata/clockface'

// Hooks
import {useDispatch} from 'react-redux'

// Actions
import {setStep} from 'src/setup/actions/creators'

// Types
import {Step} from 'src/setup/reducers'

export const Welcome: FC = () => {
  const dispatch = useDispatch()

  const handleNext = () => {
    dispatch(setStep(Step.Admin))
  }

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
        testID={'get-start'}
        onClick={handleNext}
      />
    </div>
  )
}

export default Welcome
