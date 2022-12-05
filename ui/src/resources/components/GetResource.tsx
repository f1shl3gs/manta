import React, {FunctionComponent, useEffect} from 'react'
import {connect, ConnectedProps} from 'react-redux'

// Components
import {SpinnerContainer, TechnoSpinner} from '@influxdata/clockface'

// Types
import {AppState} from 'src/types/stores'
import {Resource, ResourceType} from 'src/types/resources'

// Actions
import {getDashboard} from 'src/dashboards/actions/thunks'
import {getCell} from 'src/cells/actions/thunk'

// Selectors
import {getResourceStatus} from 'src/resources/selectors'

interface OwnProps {
  resources: Resource[]
  children: React.ReactNode
}

const mstp = (state: AppState, props: OwnProps) => {
  const loading = getResourceStatus(state, props.resources)

  return {
    loading,
  }
}

const mdtp = {
  getCell,
  getDashboard,
}

const connector = connect(mstp, mdtp)

type Props = OwnProps & ConnectedProps<typeof connector>

const GetResource: FunctionComponent<Props> = props => {
  const {loading, children, resources, getDashboard, getCell} = props

  useEffect(() => {
    const getResourceDetails = ({type, id}: Resource) => {
      switch (type) {
        case ResourceType.Dashboards:
          return getDashboard(id)

        case ResourceType.Cells:
          return getCell(id)

        default:
          throw new Error(
            `incorrect resouce type: ${type} provided to GetResource`
          )
      }
    }

    const promises = []

    resources.forEach(resource => {
      promises.push(getResourceDetails(resource))
    })

    Promise.all(promises)
  }, [resources, getCell, getDashboard])

  return (
    <SpinnerContainer loading={loading} spinnerComponent={<TechnoSpinner />}>
      {children}
    </SpinnerContainer>
  )
}

export default connector(GetResource)
