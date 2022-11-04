import React, {FunctionComponent} from 'react'
import {
  RemoteDataState,
  SpinnerContainer,
  TechnoSpinner,
} from '@influxdata/clockface'

interface Props {
  children?: JSX.Element | JSX.Element[]
  loading?: RemoteDataState
}

const PageSpinner: FunctionComponent<Props> = ({children, loading}) => {
  return (
    <SpinnerContainer
      loading={loading ?? RemoteDataState.Loading}
      spinnerComponent={<TechnoSpinner />}
    >
      {children}
    </SpinnerContainer>
  )
}

export default PageSpinner
