// Libraries
import React, {useState} from 'react'
import {useHistory} from 'react-router-dom'

// Components
import {ComponentSize, Orientation, Tabs} from '@influxdata/clockface'

interface Props {
  tabs: {
    id: string
    text: string
  }[]
  prefix: string
}

const AlertsNavigation: React.FC<Props> = props => {
  const {prefix, tabs} = props
  const [activeTab, setActive] = useState(() => {
    if (window.location.pathname.indexOf('checks') >= 0) {
      return 'checks'
    }

    return 'endpoints'
  })

  const history = useHistory()

  const onClick = (id: string) => {
    setActive(id)
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
          key={tab.id}
          text={tab.text}
          id={tab.id}
          // @ts-ignore
          onClick={onClick}
          active={tab.id === activeTab}
        />
      ))}
    </Tabs>
  )
}

export default AlertsNavigation
