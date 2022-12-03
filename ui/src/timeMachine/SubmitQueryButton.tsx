import React from 'react'
import {Button, ComponentColor} from '@influxdata/clockface'
import {useAutoRefresh} from 'src/shared/useAutoRefresh'

const SubmitQueryButton: React.FC = () => {
  const {refresh} = useAutoRefresh()

  return (
    <Button text={'Submit'} color={ComponentColor.Primary} onClick={refresh} />
  )
}

export default SubmitQueryButton
