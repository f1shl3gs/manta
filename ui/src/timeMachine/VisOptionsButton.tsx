import React from 'react'

import {Button, ComponentColor, IconFont} from '@influxdata/clockface'
import {useViewingVisOptions} from 'src/timeMachine/useTimeMachine'

const VisOptionsButton: React.FC = () => {
  const {viewingVisOptions, setViewingVisOptions} = useViewingVisOptions()

  const onToggleVisOptions = () => {
    setViewingVisOptions(prev => !prev)
  }

  const color = viewingVisOptions
    ? ComponentColor.Primary
    : ComponentColor.Default

  return (
    <Button
      color={color}
      icon={IconFont.CogSolid_New}
      onClick={onToggleVisOptions}
      text={'Customize'}
    />
  )
}

export default VisOptionsButton
