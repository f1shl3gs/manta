// Libraries
import React, {useRef, useState} from 'react'
import classnames from 'classnames'

interface Props {}

const FlameGraphPanel: React.FC<Props> = props => {
  const [view, setView] = useState('icicle')
  const [viewType, setViewType] = useState('double')

  const tablePane = classnames('pane', {
    hidden: view === 'icicle',
    verticalOrientation: viewType === 'double',
  })

  return <canvas className={'flamegraph-canvas'} height={0} />
}

export default FlameGraphPanel
