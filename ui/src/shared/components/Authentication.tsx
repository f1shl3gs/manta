import React, {FunctionComponent} from 'react'
import {useAuthentication} from 'src/shared/components/useAuthentication'
import PageSpinner from 'src/shared/components/PageSpinner'

interface Props {
  children: any
}

const Authentication: FunctionComponent<Props> = ({children}) => {
  const {loading} = useAuthentication()

  return <PageSpinner loading={loading}>{children}</PageSpinner>
}

export default Authentication
