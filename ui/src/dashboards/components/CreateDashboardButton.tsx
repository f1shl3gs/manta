// Libraries
import React, {FunctionComponent} from 'react'
import {useDispatch} from 'react-redux'
import {useNavigate} from 'react-router-dom'

// Components
import AddResourceDropdown from 'src/shared/components/AddResourceDropdown'

// Actions
import {createDashboard} from 'src/dashboards/actions/thunks'

const CreateDashboardButton: FunctionComponent = () => {
  const dispatch = useDispatch()
  const navigate = useNavigate()

  const onSelectNew = (): void => {
    dispatch(createDashboard())
  }

  const onSelectImport = (): void => {
    navigate(`${window.location.pathname}/import`)
  }

  return (
    <AddResourceDropdown
      resourceType={'dashboard'}
      onSelectNew={onSelectNew}
      onSelectImport={onSelectImport}
    />
  )
}

export default CreateDashboardButton
