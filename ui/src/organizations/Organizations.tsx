// Libraries
import React, {FunctionComponent, ReactNode} from 'react'

// Components
import {OrganizationsProvider} from 'src/organizations/useOrganizations'
import Nav from 'src/layout/Nav'
import PageSpinner from 'src/shared/components/PageSpinner'

// Types
import {Organization} from 'src/types/organization'

// Hooks
import useFetch from 'src/shared/useFetch'

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
