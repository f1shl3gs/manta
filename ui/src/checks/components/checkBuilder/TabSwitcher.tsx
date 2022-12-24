// Libraries
import React, {FunctionComponent} from 'react'
import {connect, ConnectedProps} from 'react-redux'

// Components
import {ButtonShape, ComponentSize, SelectGroup} from '@influxdata/clockface'

// Types
import {AppState} from 'src/types/stores'

// Actions
import {setTab} from 'src/checks/actions/builder'

const mstp = (state: AppState) => {
  const {tab} = state.checkBuilder

  return {
    tab,
  }
}

const mdtp = {
  setTab,
}

const connector = connect(mstp, mdtp)
type Props = ConnectedProps<typeof connector>

const TabSwitcher: FunctionComponent<Props> = ({tab, setTab}) => {
  return (
    <SelectGroup
      shape={ButtonShape.StretchToFit}
      style={{width: '400px'}}
      size={ComponentSize.Small}
    >
      <SelectGroup.Option
        id={'query'}
        active={tab === 'query'}
        value={'query'}
        onClick={setTab}
      >
        1. Define Query
      </SelectGroup.Option>

      <SelectGroup.Option
        id={'meta'}
        active={tab === 'meta'}
        value={'meta'}
        onClick={setTab}
      >
        2. Configure check
      </SelectGroup.Option>
    </SelectGroup>
  )
}

export default connector(TabSwitcher)
