import React from 'react'

import {Button, ComponentColor, IconFont} from '@influxdata/clockface'
import {useViewOption} from 'src/dashboards/components/useViewOption'

const VisOptionsButton: React.FC = () => {
  const {isViewingVisOptions, onToggleVisOptions} = useViewOption()

  const color = isViewingVisOptions
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
