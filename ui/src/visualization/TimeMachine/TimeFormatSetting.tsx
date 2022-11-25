import React from 'react'
import {SelectDropdown} from '@influxdata/clockface'
import {FORMAT_OPTIONS, resolveTimeFormat} from 'src/constants/timeFormat'

interface Props {
  timeFormat: string
  onTimeFormatChange: (format: string) => void
}

const TimeFormatSetting: React.FC<Props> = props => {
  const {timeFormat, onTimeFormatChange} = props

  return (
    <SelectDropdown
      options={FORMAT_OPTIONS.map(option => option.text)}
      selectedOption={resolveTimeFormat(timeFormat)}
      onSelect={onTimeFormatChange}
    />
  )
}

export default TimeFormatSetting
