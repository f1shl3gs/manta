// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {StepSwitcher} from './components/StepSwitcher'
import {AppWrapper} from '@influxdata/clockface'
import {WizardFullScreen} from '../shared/components/wizard/WizardFullScreen'
import {StepProvider} from './useStep'

const Setup: FunctionComponent = () => {
  return (
    <AppWrapper>
      <StepProvider>
        <WizardFullScreen>
          <div className={'wizard-contents'}>
            <div className={'wizard-step--container'}>
              <StepSwitcher />
            </div>
          </div>
        </WizardFullScreen>
      </StepProvider>
    </AppWrapper>
  )
}

export default Setup
