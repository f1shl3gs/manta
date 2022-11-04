import React, {FunctionComponent} from 'react'
import {useAuthentication} from './useAuthentication'
import PageSpinner from './PageSpinner'

interface Props {
  children: any
}

const Authentication: FunctionComponent<Props> = ({children}) => {
  const {loading} = useAuthentication()

  return <PageSpinner loading={loading}>{children}</PageSpinner>
}

export default Authentication
