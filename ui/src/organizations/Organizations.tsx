import React, {FunctionComponent, ReactNode} from 'react'
import useFetch from 'src/shared/useFetch'
import {Organization} from 'src/types/organization'
import {OrganizationsProvider} from 'src/organizations/useOrganizations'
import Nav from 'src/organizations/Nav'
import PageSpinner from 'src/shared/components/PageSpinner'

interface Props {
  children: ReactNode
}

const Organizations: FunctionComponent<Props> = ({children}) => {
  const {
    data = [],
    loading,
    run: refetch,
  } = useFetch<[Organization]>('/api/v1/organizations')

  return (
    <PageSpinner loading={loading}>
      <OrganizationsProvider organizations={data} refetch={refetch}>
        <Nav />

        {children}
      </OrganizationsProvider>
    </PageSpinner>
  )
}

export default Organizations
