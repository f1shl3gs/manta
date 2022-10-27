import React, {FunctionComponent, ReactNode, useState} from 'react'
import {SpinnerContainer, TechnoSpinner} from '@influxdata/clockface'
import useFetch from '../useFetch'
import {useParams} from 'react-router-dom'
import constate from 'constate'

enum ResourceType {
  Dashboards = 'dashboards',
  Users = 'users',
}

interface State<T> {
  resources: T[]
  reload: () => void
}

const [ResourcesProvider, useResources] = constate((state: State<any>) => state)

interface Props {
  children: ReactNode

  type: ResourceType
  url?: string
}

const GetResources: FunctionComponent<Props> = ({children, type, url}) => {
  const {orgId} = useParams()
  const u = url ? url : `/api/v1/${type}?orgId=${orgId}`
  const {run, data, loading} = useFetch(u)


  return (
    <SpinnerContainer spinnerComponent={<TechnoSpinner />} loading={loading}>
      <ResourcesProvider
        resources={data}
        reload={run}
      >
        {children}
      </ResourcesProvider>
    </SpinnerContainer>
  )
}

export {ResourceType, GetResources, useResources}
