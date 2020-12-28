import React, {useState} from 'react'
import {ClickOutside, ComponentSize, Input} from '@influxdata/clockface'

interface Props {
  onUpdate: (val: string) => void
  value?: string
  onClick?: (e: MouseEvent) => void
  placeholder?: string
}

const Editable: React.FC<Props> = ({placeholder, value}) => {
  const [editing, setEditing] = useState(false)
  const [text, setText] = useState(value)

  return (
    <ClickOutside onClickOutside={() => setEditing(false)}>
      <Input
        className="cf-resource-editable-name--input"
        size={ComponentSize.Medium}
        maxLength={90}
        placeholder={placeholder}
        autoFocus={false}
        value={text}
        onChange={(e) => {
          setText(e.target.value)
        }}
      />
    </ClickOutside>
  )
}

export default Editable