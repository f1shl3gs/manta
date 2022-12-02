import React, {FunctionComponent} from 'react'
import {ViewProperties} from 'src/types/dashboard'
import TimeSeries from '../../../shared/components/TimeSeries'

interface Props {
  id: string
  properties: ViewProperties
}

const RefreshingView: FunctionComponent<Props> = ({properties}) => {
  return <TimeSeries viewProperties={properties} />
}

export default RefreshingView
