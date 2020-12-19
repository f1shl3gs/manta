// Libraries
import React, {PureComponent} from 'react'

interface Props {
  id: string
  manualRefresh: number
}

const RefreshingView: React.FC<Props> = props => {
  return (
    <div>RefreshView</div>
  )
}

export default RefreshingView;

