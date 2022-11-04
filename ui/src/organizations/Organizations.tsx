import React, {FunctionComponent, ReactNode} from 'react'
import useFetch from 'shared/useFetch'
import {Organization} from 'types/Organization'
import {OrganizationsProvider} from './useOrganizations'
import Nav from './Nav'
import PageSpinner from '../shared/components/PageSpinner';

interface Props {
  children: ReactNode
}

const Organizations: FunctionComponent<Props> = ({children}) => {
  const {data, loading} = useFetch<[Organization]>('/api/v1/organizations')

  return (
    <PageSpinner loading={loading}>
      <OrganizationsProvider organizations={data || []}>
        <Nav />

        {children}
      </OrganizationsProvider>
    </PageSpinner>
  )
}

export default Organizations
