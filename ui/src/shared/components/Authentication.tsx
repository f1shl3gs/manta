import React, {FunctionComponent} from 'react'
import {SpinnerContainer, TechnoSpinner} from '@influxdata/clockface'
import {useAuthentication} from './useAuthentication'

interface Props {
  children: any
}

const Authentication: FunctionComponent<Props> = ({children}) => {
  const {loading} = useAuthentication()

  return (
    <SpinnerContainer loading={loading} spinnerComponent={<TechnoSpinner />}>
      {children}
    </SpinnerContainer>
  )
}

export default Authentication
