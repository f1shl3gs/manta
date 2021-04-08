import React from 'react'

interface Props {
  prefix: string
  tabs: {
    id: string
    text: string
  }[]
}

const SettingsNavigation: React.FC<Props> = props => {
  const {prefix} = props
  console.log('prefix', prefix)
  return <div>aaaa</div>
}

export default SettingsNavigation
