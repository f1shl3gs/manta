import React from 'react'
import classnames from 'classnames'

interface Props {
  id: string
  name: string
  active: boolean
  onClick: (id: string) => void

  testID?: string
}

const ToolbarTab: React.FC<Props> = props => {
  const {id, name, active, testID = 'toolbar-tab', onClick} = props
  const toolbarTabClass = classnames('flux-toolbar--tab', {
    'flux-toolbar--tab__active': active,
  })

  const handleClick = () => {
    if (active) {
      onClick('none')
    } else {
      onClick(id)
    }
  }

  return (
    <div
      className={toolbarTabClass}
      onClick={handleClick}
      title={name}
      data-testid={testID}
    >
      <div className={'flux-toolbar--tab-label'}>{name}</div>
    </div>
  )
}

export default ToolbarTab
