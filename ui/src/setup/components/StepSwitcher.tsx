// Library
import React, {FC} from 'react'

// Components
import Admin from 'src/setup/components/Admin'
import Welcome from 'src/setup/components/Welcome'

// Hooks
import {useSelector} from 'react-redux'

// Types
import {AppState} from 'src/types/stores'
import {Step} from 'src/setup/reducers'

export const StepSwitcher: FC = () => {
  const step = useSelector((state: AppState) => state.setup.step)

  switch (step) {
    case Step.Welcome:
      return <Welcome />
    case Step.Admin:
      return <Admin />
  }
}
