import React from 'react'
import {InsideThreshold, OutsideThreshold} from '../../../types/Check'
import {convertUserInputToNumOrNaN} from '../../../utils/convert'
import {FlexBox, Input, InputType, TextBlock} from '@influxdata/clockface'

interface Props {
  threshold: InsideThreshold | OutsideThreshold
  changeRange: (min: number, max: number) => void
}

const ThresholdRangeInput: React.FC<Props> = props => {
  const {threshold, changeRange} = props

  const onChangeMin = (e: React.ChangeEvent<HTMLInputElement>) => {
    const min = convertUserInputToNumOrNaN(e)

    changeRange(min, threshold.max)
  }

  const onChangeMax = (e: React.ChangeEvent<HTMLInputElement>) => {
    const max = convertUserInputToNumOrNaN(e)

    changeRange(threshold.min, max)
  }

  return (
    <>
      <FlexBox.Child testID={'component-spacer--flex-child'}>
        <Input
          onChange={onChangeMin}
          name={'min'}
          testID={'input-field'}
          type={InputType.Number}
          value={threshold.min}
        />
      </FlexBox.Child>
      <TextBlock testID={'text-block'} text={'to'} />
      <FlexBox.Child testID={'component-spacer--flex-child'}>
        <Input
          onChange={onChangeMax}
          name={'max'}
          testID={'input-field'}
          type={InputType.Number}
          value={threshold.max}
        />
      </FlexBox.Child>
    </>
  )
}

export default ThresholdRangeInput
