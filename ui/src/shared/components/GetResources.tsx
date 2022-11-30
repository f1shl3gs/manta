// Libraries
import React, {FunctionComponent, ReactNode} from 'react'
import {useParams} from 'react-router-dom'
import constate from 'constate'

// Components
import PageSpinner from 'src/shared/components/PageSpinner'

// Hooks
import useFetch from 'src/shared/useFetch'

// Types
import {ResourceType} from 'src/types/resources'

interface State<T> {
  resources: T[]
  reload: () => void
}

const [ResourcesProvider, useResources] = constate((state: State<any>) => state)

interface Props {
  type: ResourceType
  url?: string

  children: ReactNode
}

const GetResources: FunctionComponent<Props> = ({type, url, children}) => {
  const {orgID} = useParams()
  const u = url ? url : `/api/v1/${type}?orgID=${orgID}`
  const {run, data, loading} = useFetch(u)

  return (
    <PageSpinner loading={loading}>
      <ResourcesProvider resources={data} reload={run}>
        {children}
      </ResourcesProvider>
    </PageSpinner>
  )
}

const withResources = (
  WrappedComponent: FunctionComponent,
  type: ResourceType
) => {
  const component = () => {
    return (
      <GetResources type={type}>
        <WrappedComponent />
      </GetResources>
    )
  }

  component.displayName = `withResources(${WrappedComponent.displayName})`

  return component
}

export {ResourceType, GetResources, useResources, withResources}
