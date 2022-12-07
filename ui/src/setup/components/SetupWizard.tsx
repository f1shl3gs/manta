// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {StepSwitcher} from 'src/setup/components/StepSwitcher'
import {AppWrapper} from '@influxdata/clockface'
import {WizardFullScreen} from 'src/shared/components/wizard/WizardFullScreen'

const Setup: FunctionComponent = () => {
  return (
    <AppWrapper>
      <WizardFullScreen>
        <div className={'wizard-contents'}>
          <div className={'wizard-step--container'}>
            <StepSwitcher />
          </div>
        </div>
      </WizardFullScreen>
    </AppWrapper>
  )
}

export default Setup
