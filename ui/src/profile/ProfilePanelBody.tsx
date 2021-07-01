// Libraries
import React from 'react'
import {useViewType, ViewType} from './useViewType'
import TablePanel from './components/TablePanel'
import FlameGraphPanel from './components/FlameGraphPanel'

const ProfilePanelBody: React.FC = () => {
  const {viewType} = useViewType()
  let content = null

  if (viewType === ViewType.Both) {
    content = (
      <>
        <TablePanel />
        <FlameGraphPanel />
      </>
    )
  } else if (viewType === ViewType.FlameGraph) {
    content = <FlameGraphPanel />
  } else {
    content = <TablePanel />
  }

  return content
}

export default ProfilePanelBody
