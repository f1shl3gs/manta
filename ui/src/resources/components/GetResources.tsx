// Libraries
import React, {FunctionComponent, useEffect} from 'react'
import {connect, ConnectedProps} from 'react-redux'

import {ResourceType} from 'src/types/resources'
import PageSpinner from 'src/shared/components/PageSpinner'
import {AppState} from 'src/types/stores'
import {getResourcesStatus} from 'src/resources/selectors'
import { getDashboards } from 'src/dashboards/actions/thunks'
import {RemoteDataState} from '@influxdata/clockface'

type ReduxProps = ConnectedProps<typeof connector>
interface OwnProps {
  loading: RemoteDataState
  resources: Array<ResourceType>
  childrent: JSX.Element | JSX.Element[]
}

type Props = ReduxProps & OwnProps

const getResourceDetails = (resource: ResourceType, props: ReduxProps) => {
  switch (resource) {
    case ResourceType.Dashboards:
      return props.getDashboards()
    default:
      throw new Error('incorrent resource type provided')
  }
}

const GetResources: FunctionComponent<Props> = (props) => {
  const {resources, loading, childrent} = props

  useEffect(() => {
    const promises = []

    resources.forEach(resource => {
      promises.push(getResourceDetails(resource, props))
    })

    Promise.all(promises)
  }, [resources, props])

  return (
    <PageSpinner loading={loading}>
      {childrent}
    </PageSpinner>
  )
}

const mstp = (state: AppState, {resources}: OwnProps) => {
  const loading = getResourcesStatus(state, resources)

  return {
    loading
  }
}

const mdtp = {
  getDashboards: getDashboards
}

const connector = connect(mstp, mdtp)

export default connector(GetResources)
