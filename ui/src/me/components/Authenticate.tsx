import React, {FunctionComponent, useEffect} from 'react'
import {connect, ConnectedProps} from 'react-redux'
import PageSpinner from 'src/shared/components/PageSpinner'
import {AppState} from 'src/types/stores'
import {getMe} from 'src/me/actions/thunks'

interface OwnProps {
  children: JSX.Element | JSX.Element[]
}

const mstp = (state: AppState) => {
  return {
    loading: state.me.state,
  }
}

const mdtp = {
  getMe,
}

const connector = connect(mstp, mdtp)

type Props = OwnProps & ConnectedProps<typeof connector>

export const Authenticate: FunctionComponent<Props> = ({
  children,
  loading,
  getMe,
}) => {
  useEffect(() => {
    getMe()
  }, [getMe])

  return <PageSpinner loading={loading}>{children}</PageSpinner>
}

export default connector(Authenticate)
