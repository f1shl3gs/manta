// Libraries
import React, {FunctionComponent, useEffect} from 'react'
import {connect, ConnectedProps} from 'react-redux'

import {ResourceType} from 'src/types/resources'
import PageSpinner from 'src/shared/components/PageSpinner'
import {AppState} from 'src/types/stores'
import {getResourcesStatus} from 'src/resources/selectors'

// Actions
import {getDashboards} from 'src/dashboards/actions/thunks'
import {getScrapes} from 'src/scrapes/actions/thunk'

interface OwnProps {
  resources: Array<ResourceType>
  children: JSX.Element | JSX.Element[]
}

type ReduxProps = ConnectedProps<typeof connector>
type Props = ReduxProps & OwnProps

const getResourceDetails = (resource: ResourceType, props: ReduxProps) => {
  switch (resource) {
    case ResourceType.Dashboards:
      return props.getDashboards()

    case ResourceType.Scrapes:
      return props.getScrapes()

    default:
      throw new Error('incorrent resource type provided')
  }
}

const GetResources: FunctionComponent<Props> = props => {
  const {resources, loading, children} = props

  useEffect(() => {
    const promises = []

    resources.forEach(resource => {
      promises.push(getResourceDetails(resource, props))
    })

    Promise.all(promises)
  }, [resources])

  return <PageSpinner loading={loading}>{children}</PageSpinner>
}

const mstp = (state: AppState, {resources}: OwnProps) => {
  const loading = getResourcesStatus(state, resources)

  return {
    loading,
  }
}

const mdtp = {
  getDashboards,
  getScrapes,
}

const connector = connect(mstp, mdtp)

export default connector(GetResources)
