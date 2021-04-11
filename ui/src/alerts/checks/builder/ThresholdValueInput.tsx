// Libraries
import React from 'react'

// Components
import {FlexBox, Input, InputType} from '@influxdata/clockface'

// Types
import {GreatThanThreshold, LessThanThreshold} from 'types/Check'

// Utils
import {convertUserInputToNumOrNaN} from 'utils/convert'

interface Props {
  threshold: GreatThanThreshold | LessThanThreshold
  changeValue: (v: number) => void
}

const ThresholdValueInput: React.FC<Props> = props => {
  const {threshold, changeValue} = props
  const onChangeValue = (e: React.ChangeEvent<HTMLInputElement>) => {
    changeValue(convertUserInputToNumOrNaN(e))
  }

  return (
    <FlexBox.Child testID={'component-spacer--flex-child'}>
      <Input
        onChange={onChangeValue}
        name={''}
        testID={'input-field'}
        type={InputType.Number}
        value={threshold.value}
      />
    </FlexBox.Child>
  )
}

export default ThresholdValueInput
