import React, {FunctionComponent} from 'react'
import {useOrganization} from 'src/organizations/useOrganizations'
import ImportOverlay from 'src/shared/components/ImportOverlay'
import {defaultErrorNotification} from 'src/shared/components/notifications/defaults'
import useFetch from 'src/shared/useFetch'
import {useNavigate} from 'react-router-dom'
import {useNotify} from '../shared/components/notifications/useNotification'

const DashboardImportOverlay: FunctionComponent = () => {
  const {id: orgId} = useOrganization()
  const navigate = useNavigate()
  const notify = useNotify()
  const {run: create} = useFetch(`/api/v1/dashboards`, {
    method: 'POST',
    body: {
      orgId,
      cells: [],
    },
    onError: err => {
      notify({
        ...defaultErrorNotification,
        message: `Create new dashboard failed\n${err}`,
      })
    },
    onSuccess: dashboard => {
      navigate(`/orgs/${orgId}/dashboards/${dashboard.id}`)
    },
  })

  const onSubmit = (imported: string) => {
    const dashboard = JSON.parse(imported)
    create(dashboard)
  }

  return <ImportOverlay resourceName={'Dashboard'} onSubmit={onSubmit} />
}

export default DashboardImportOverlay
