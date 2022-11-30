import React, {FunctionComponent} from 'react'
import ImportOverlay from 'src/shared/components/ImportOverlay'
import {defaultErrorNotification} from 'src/shared/components/notifications/defaults'
import useFetch from 'src/shared/useFetch'
import {useNavigate} from 'react-router-dom'
import {useNotify} from 'src/shared/components/notifications/useNotification'
import {useOrg} from 'src/organizations/selectors'

const DashboardImportOverlay: FunctionComponent = () => {
  const {id: orgID} = useOrg()
  const navigate = useNavigate()
  const notify = useNotify()
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
      navigate(`/orgs/${orgID}/dashboards/${dashboard.id}`)
    },
  })

  const onSubmit = (imported: string) => {
    const dashboard = JSON.parse(imported)
    create(dashboard)
  }

  return <ImportOverlay resourceName={'Dashboard'} onSubmit={onSubmit} />
}

export default DashboardImportOverlay
