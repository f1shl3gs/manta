import React from 'react'

interface Props {
  query: string
}

const CheckVis: React.FC<Props> = ({query}) => {
  return <div>todo: {query}</div>
}

export default CheckVis
