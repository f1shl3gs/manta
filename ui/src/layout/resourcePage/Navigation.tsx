// Libraries
import React, {useState} from 'react'
import {useHistory} from 'react-router-dom'

// Components
import {ComponentSize, Orientation, Tabs} from '@influxdata/clockface'

interface Props {
  prefix: string
  tabs: {
    id: string
    text: string
  }[]
}

const Navigation: React.FC<Props> = props => {
  const {prefix, tabs} = props
  const [activeTab, setActiveTab] = useState(() => {
    for (let tab of tabs) {
      if (window.location.pathname.indexOf(tab.id) >= 0) {
        return tab.id
      }
    }

    return tabs[0].id
  })

  const history = useHistory()

  const onClick = (id: string) => {
    setActiveTab(id)
    history.push(`${prefix}/${id}`)
  }

  return (
    <Tabs
      orientation={Orientation.Horizontal}
      size={ComponentSize.Large}
      dropdownBreakpoint={872}
      dropdownLabel={''}
    >
      {tabs.map(tab => (
        <Tabs.Tab
          id={tab.id}
          key={tab.id}
          text={tab.text}
          // @ts-ignore
          onClick={onClick}
          active={tab.id === activeTab}
        />
      ))}
    </Tabs>
  )
}

export default Navigation
