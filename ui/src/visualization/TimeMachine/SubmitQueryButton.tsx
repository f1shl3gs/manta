import React from 'react'
import {Button, ComponentColor} from '@influxdata/clockface'
import {useQueries} from './useQueries'

const SubmitQueryButton: React.FC = () => {
  const {activeIndex, activeQuery} = useQueries()
  console.log(activeIndex, activeQuery)

  return (
    <Button
      text={'Submit'}
      color={ComponentColor.Primary}
      onClick={() => console.log('submit')}
    />
  )
}

export default SubmitQueryButton
