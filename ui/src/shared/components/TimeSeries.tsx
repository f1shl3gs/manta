import React, {FunctionComponent} from 'react'
import {ViewProperties} from 'src/types/Dashboard'

interface Props {
  cellID?: string
  viewProperties: ViewProperties
}

const TimeSeries: FunctionComponent<Props> = props => {
  const {viewProperties} = props

  return <></>
}

export default TimeSeries
