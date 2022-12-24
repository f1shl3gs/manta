// Libraries
import React, {FunctionComponent} from 'react'
import {connect, ConnectedProps} from 'react-redux'

// Components
import QueryEditor from 'src/timeMachine/components/queryEditor/QueryEditor'

// Types
import {AppState} from 'src/types/stores'

// Actions
import {setQuery} from 'src/checks/actions/builder'

const mstp = (state: AppState) => {
  const {query} = state.checkBuilder

  return {
    query,
  }
}

const mdtp = {
  setQuery,
}

const connector = connect(mstp, mdtp)
type Props = ConnectedProps<typeof connector>

const CheckQueryEditor: FunctionComponent<Props> = ({query, setQuery}) => {
  return <QueryEditor query={query} onChange={setQuery} />
}

export default connector(CheckQueryEditor)
