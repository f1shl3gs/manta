import React from 'react'
import {Button, ComponentColor} from '@influxdata/clockface'

const SubmitQueryButton: React.FC = () => {
  return (
    <Button text={'Submit'} color={ComponentColor.Primary} onClick={refresh} />
  )
}

export default SubmitQueryButton
