import React, {FunctionComponent} from 'react'
import {XYViewProperties} from 'src/types/Dashboard'
import {VisualizationProps} from 'src/visualization'

interface Props extends VisualizationProps {
  properties: XYViewProperties
}

const Line: FunctionComponent<Props> = props => {
  return <></>
}

export default Line
