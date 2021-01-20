import React from 'react'

import {Button, ComponentColor, IconFont} from '@influxdata/clockface'
import {useViewOption} from '../../shared/useViewOption'

const VisOptionsButton: React.FC = () => {
  const {isViewingVisOptions, onToggleVisOptions} = useViewOption()

  const color = isViewingVisOptions
    ? ComponentColor.Primary
    : ComponentColor.Default

  return (
    <Button
      color={color}
      icon={IconFont.CogThick}
      onClick={onToggleVisOptions}
      text={'Customize'}
    />
  )
}

export default VisOptionsButton
