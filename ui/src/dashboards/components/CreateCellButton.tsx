import React, {FunctionComponent, MouseEvent} from 'react'

import {Button, ComponentColor, IconFont} from '@influxdata/clockface'
import {useDashboard} from 'src/dashboards/useDashboard'

const CreateCellButton: FunctionComponent = () => {
  const {addCell} = useDashboard()

  const handleAddCell = (_ev?: MouseEvent<HTMLButtonElement>): void => {
    addCell()
  }

  return (
    <Button
      text={'Add Cell'}
      color={ComponentColor.Primary}
      icon={IconFont.AddCell_New}
      onClick={handleAddCell}
    />
  )
}

export default CreateCellButton
