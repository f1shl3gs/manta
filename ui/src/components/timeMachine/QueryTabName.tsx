// Libraries
import React from 'react'

interface Props {
  name: string
}

const QueryTabName: React.FC<Props> = props => {
  const {name} = props

  return (
    <div className={'query-tab--name'} title={name}>
      {name}
    </div>
  )
}

export default QueryTabName
