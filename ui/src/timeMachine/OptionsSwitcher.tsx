import React from 'react'

import {ViewProperties} from 'src/types/cells'
import {LineOptions} from 'src/visualization/Line'

interface Props {
  view: ViewProperties
}

const OptionsSwitcher: React.FC<Props> = props => {
  const {view} = props

  switch (view.type) {
    case 'gauge':
      return <div>todo</div>

    case 'xy':
      return <LineOptions />
    default:
      return <div>Unknown</div>
  }
}

export default OptionsSwitcher
