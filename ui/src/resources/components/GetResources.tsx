// Libraries
import React, {FunctionComponent, useEffect} from 'react'
import {connect, ConnectedProps} from 'react-redux'

// Components
import PageSpinner from 'src/shared/components/PageSpinner'

// Types
import {ResourceType} from 'src/types/resources'
import {AppState} from 'src/types/stores'

// Actions
import {getDashboards} from 'src/dashboards/actions/thunks'
import {getScrapes} from 'src/scrapes/actions/thunk'
import {getMembers} from 'src/members/actions/thunk'
import {getConfigs} from 'src/configurations/actions/thunk'

// Selectors
import {getResourcesStatus} from 'src/resources/selectors'

const mstp = (state: AppState, {resources}: OwnProps) => {
  const loading = getResourcesStatus(state, resources)

  return {
    loading,
  }
}

const mdtp = {
  getConfigs,
  getDashboards,
  getMembers,
  getScrapes,
}

const connector = connect(mstp, mdtp)

interface OwnProps {
  resources: Array<ResourceType>
  children: JSX.Element | JSX.Element[]
}

type ReduxProps = ConnectedProps<typeof connector>
type Props = ReduxProps & OwnProps

const GetResources: FunctionComponent<Props> = props => {
  const {
    resources,
    loading,
    children,
    getConfigs,
    getDashboards,
    getMembers,
    getScrapes,
  } = props

  useEffect(() => {
    const getResourceDetails = (resource: ResourceType) => {
      switch (resource) {
        case ResourceType.Configurations:
          return getConfigs()

        case ResourceType.Dashboards:
          return getDashboards()

        case ResourceType.Members:
          return getMembers()

        case ResourceType.Scrapes:
          return getScrapes()

        default:
          throw new Error('incorrent resource type provided')
      }
    }

    const promises = []

    resources.forEach(resource => {
      promises.push(getResourceDetails(resource))
    })

    Promise.all(promises)
  }, [resources, getConfigs, getDashboards, getMembers, getScrapes])

  return <PageSpinner loading={loading}>{children}</PageSpinner>
}

export default connector(GetResources)
