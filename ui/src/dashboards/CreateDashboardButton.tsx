// Libraries
import React, {FunctionComponent} from 'react'

// Components
import AddResourceDropdown from '../shared/components/AddResourceDropdown'

// Hooks
import {
  defaultErrorNotification,
  useNotify,
} from 'src/shared/components/notifications/useNotification'
import {useNavigate} from 'react-router-dom'
import useFetch from 'src/shared/useFetch'
import {useOrg} from 'src/organizations/selectors'

const CreateDashboardButton: FunctionComponent = () => {
  const {id: orgID} = useOrg()
  const notify = useNotify()
  const navigate = useNavigate()
  const {run: create} = useFetch(`/api/v1/dashboards`, {
    method: 'POST',
    body: {
      orgID: orgID,
      cells: [],
    },
    onError: err => {
      notify({
        ...defaultErrorNotification,
        message: `Create new dashboard failed\n${err}`,
      })
    },
    onSuccess: dashboard => {
      navigate(`${window.location.pathname}/${dashboard.id}`)
    },
  })

  const onSelectNew = (): void => {
    create({
      name: '',
    })
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
