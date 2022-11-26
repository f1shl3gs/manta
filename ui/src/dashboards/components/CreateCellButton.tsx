// Libraries
import React, {FunctionComponent, MouseEvent} from 'react'
import {useNavigate} from 'react-router-dom'

// Components
import {Button, ComponentColor, IconFont} from '@influxdata/clockface'

const CreateCellButton: FunctionComponent = () => {
  const navigate = useNavigate()

  const handleAddCell = (_ev?: MouseEvent<HTMLButtonElement>): void => {
    navigate(`${window.location.pathname}/cells/new`)
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
