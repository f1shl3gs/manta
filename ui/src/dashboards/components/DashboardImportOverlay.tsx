// Libraries
import {connect, ConnectedProps} from 'react-redux'
import React, {FunctionComponent} from 'react'

// Components
import ImportOverlay from 'src/shared/components/ImportOverlay'

// Actions
import {createDashboardFromJSON} from 'src/dashboards/actions/thunks'

const mdtp = {
  create: createDashboardFromJSON,
}

const connector = connect(null, mdtp)

type Props = ConnectedProps<typeof connector>

const DashboardImportOverlay: FunctionComponent<Props> = ({create}) => {
  const onSubmit = (imported: string) => {
    const dashboard = JSON.parse(imported)
    create(dashboard)
  }

  return <ImportOverlay resourceName={'Dashboard'} onSubmit={onSubmit} />
}

export default connector(DashboardImportOverlay)
