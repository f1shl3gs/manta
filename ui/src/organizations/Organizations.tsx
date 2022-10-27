import React, {FunctionComponent, ReactNode} from 'react'
import {SpinnerContainer, TechnoSpinner} from '@influxdata/clockface'
import useFetch from 'shared/useFetch'
import {Organization} from 'types/Organization'
import {OrganizationsProvider} from './useOrganizations'
import Nav from './Nav'

interface Props {
  children: ReactNode
}

const Organizations: FunctionComponent<Props> = ({children}) => {
  const {data, loading} = useFetch<[Organization]>('/api/v1/organizations')

  return (
    <SpinnerContainer loading={loading} spinnerComponent={<TechnoSpinner />}>
      <OrganizationsProvider organizations={data || []}>
        <Nav />

        {children}
      </OrganizationsProvider>
    </SpinnerContainer>
  )
}

export default Organizations
