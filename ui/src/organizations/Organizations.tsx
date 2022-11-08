import React, {FunctionComponent, ReactNode} from 'react'
import useFetch from 'src/shared/useFetch'
import {Organization} from 'src/types/Organization'
import {OrganizationsProvider} from 'src/organizations/useOrganizations'
import Nav from 'src/organizations/Nav'
import PageSpinner from 'src/shared/components/PageSpinner'

interface Props {
  children: ReactNode
}

const Organizations: FunctionComponent<Props> = ({children}) => {
  const {data = [], loading} = useFetch<[Organization]>('/api/v1/organizations')

  return (
    <PageSpinner loading={loading}>
      <OrganizationsProvider organizations={data}>
        <Nav />

        {children}
      </OrganizationsProvider>
    </PageSpinner>
  )
}

export default Organizations
