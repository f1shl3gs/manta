import React from 'react'

import {usePresentationMode} from 'src/shared/usePresentationMode'
import {IconFont, SquareButton} from '@influxdata/clockface'

const PresentationModeToggle = () => {
  const {togglePresentationMode} = usePresentationMode()

  return (
    <SquareButton icon={IconFont.ExpandB} onClick={togglePresentationMode} />
  )
}

export default PresentationModeToggle
