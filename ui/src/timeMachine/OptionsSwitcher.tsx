import React from 'react'

import {ViewProperties} from 'src/types/cells'
import {LineOptions} from 'src/visualization/Line'

interface Props {
  viewProperties: ViewProperties
  update: (viewProperties: any) => void
}

const OptionsSwitcher: React.FC<Props> = ({viewProperties, update}) => {
  switch (viewProperties.type) {
    case 'line-plus-single-stat':
      break
    case 'single-stat':
      break
    case 'histogram':
      break
    case 'markdown':
      break
    case 'band':
      break
    case 'table':
      break
    case 'scatter':
      break
    case 'mosaic':
      break
    case 'heatmap':
      break
    case 'gauge':
      return <div>todo</div>

    case 'xy':
      return <LineOptions viewProperties={viewProperties} update={update} />
    default:
      return <div>Unknown</div>
  }
}

export default OptionsSwitcher
