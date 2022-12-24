// Libraries
import React, {FunctionComponent, ChangeEvent} from 'react'

// Components
import {FlexBoxChild, Input, InputType} from '@influxdata/clockface'

// Types
import {InsideThreshold, OutsideThreshold, Threshold} from 'src/types/checks'

// Utils
import {convertUserInputToNumOrNaN} from 'src/shared/utils/convertUserInput'

interface Props {
  threshold: Exclude<Threshold, InsideThreshold | OutsideThreshold>
  onChange: (n: number) => void
}

const ThresholdValueInput: FunctionComponent<Props> = ({
  threshold,
  onChange,
}) => {
  const onChangeValue = (ev: ChangeEvent<HTMLInputElement>) => {
    onChange(convertUserInputToNumOrNaN(ev))
  }

  return (
    <FlexBoxChild testID={'component-spacer--flex-child'}>
      <Input
        name={''}
        onChange={onChangeValue}
        testID={'threshold-value--input'}
        type={InputType.Number}
        value={threshold.value}
      />
    </FlexBoxChild>
  )
}

export default ThresholdValueInput
