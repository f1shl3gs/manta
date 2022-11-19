import React, {FunctionComponent} from 'react'

import {Button, ComponentColor, IconFont} from '@influxdata/clockface'
import useFetch from 'src/shared/useFetch'
import {useDashboard} from 'src/dashboards/useDashboard'

const CreateCellButton: FunctionComponent = () => {
  const {id} = useDashboard()
  const {run: addCell} = useFetch(`/api/v1/dashboards/${id}/cells`, {
    method: 'POST',
    body: {
      w: 4,
      h: 4,
      x: 0,
      y: 0,
    },
  })

  return (
    <Button
      text={'Add Cell'}
      color={ComponentColor.Primary}
      icon={IconFont.AddCell_New}
      onClick={() => addCell()}
    />
  )
}

export default CreateCellButton
