// Libraries
import {useDispatch} from 'react-redux'
import React, {FunctionComponent} from 'react'

// Components
import ImportOverlay from 'src/shared/components/ImportOverlay'

// Actions
import {createDashboardFromJSON} from 'src/dashboards/actions/thunks'

const DashboardImportOverlay: FunctionComponent = () => {
  const dispatch = useDispatch()

  const onSubmit = (imported: string) => {
    const dashboard = JSON.parse(imported)
    dispatch(createDashboardFromJSON(dashboard))
  }

  return <ImportOverlay resourceName={'Dashboard'} onSubmit={onSubmit} />
}

export default DashboardImportOverlay
