import React from 'react'
import {LineOptions} from 'src/visualization/Line'
import { useTimeMachine } from './useTimeMachine'

const OptionsSwitcher: React.FC = () => {
  const {viewProperties, setViewProperties} = useTimeMachine()

  switch (viewProperties.type) {
    case 'gauge':
      return <div>todo</div>

    case 'xy':
      return (
        <LineOptions
          viewProperties={viewProperties}
          setViewProperties={setViewProperties}
        />
      )
    default:
      return <div>Unknown</div>
  }
}

export default OptionsSwitcher
