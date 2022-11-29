// Libraries
import React, {FunctionComponent, ReactNode} from 'react'

// Components
import Nav from 'src/layout/Nav'
import PageSpinner from 'src/shared/components/PageSpinner'

// Hooks
import useFetch from 'src/shared/useFetch'

// Types
import {Organization} from 'src/types/Organization'
import {OrganizationsProvider} from 'src/organizations/useOrganizations'

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
