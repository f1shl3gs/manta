// Libraries
import React, {ChangeEvent, FunctionComponent} from 'react'

// Components
import {FlexBoxChild, Input, InputType, TextBlock} from '@influxdata/clockface'

// Types
import {InsideThreshold, OutsideThreshold} from 'src/types/checks'

// Utils
import {convertUserInputToNumOrNaN} from 'src/shared/utils/convertUserInput'

interface Props {
  threshold: InsideThreshold | OutsideThreshold
  onChange: (min: number, max: number) => void
}

const ThresholdRangeInput: FunctionComponent<Props> = ({
  threshold,
  onChange,
}) => {
  const onChangeMin = (ev: ChangeEvent<HTMLInputElement>) => {
    const min = convertUserInputToNumOrNaN(ev)
    onChange(min, threshold.max)
  }

  const onChangeMax = (ev: ChangeEvent<HTMLInputElement>) => {
    const max = convertUserInputToNumOrNaN(ev)
    onChange(threshold.min, max)
  }

  return (
    <>
      <FlexBoxChild testID={'component-spacer--flex-child'}>
        <Input
          name={'min'}
          type={InputType.Number}
          value={threshold.min}
          testID={'threshold-range-min--input'}
          onChange={onChangeMin}
        />
      </FlexBoxChild>

      <TextBlock text={'to'} testID={'text-block'} />

      <FlexBoxChild>
        <Input
          name={'max'}
          type={InputType.Number}
          value={threshold.max}
          testID={'threshold-range-max--input'}
          onChange={onChangeMax}
        />
      </FlexBoxChild>
    </>
  )
}

export default ThresholdRangeInput
