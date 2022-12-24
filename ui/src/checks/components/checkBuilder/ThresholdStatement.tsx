// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {
  ButtonType,
  ComponentColor,
  ComponentSize,
  DismissButton,
  FlexBox,
  FlexBoxChild,
  FlexDirection,
  InfluxColors,
  Panel,
  PanelBody,
  SelectDropdown,
  TextBlock,
} from '@influxdata/clockface'

// Types
import {ConditionStatus, Threshold, ThresholdType} from 'src/types/checks'

// Constants
import {CHECK_STATUS_COLORS} from 'src/checks/constants'

const dropdownOptions = {
  'is above': 'gt',
  'is below': 'lt',
  'inside range': 'inside',
  'outside range': 'outside',
}

interface Props {
  status: ConditionStatus
  threshold: Threshold
  onChangeType: (to: ThresholdType) => void
  onDelete: () => void

  children: JSX.Element
}

const optionSelector = (type: ThresholdType) => {
  for (const k in dropdownOptions) {
    if (dropdownOptions[k] === type) {
      return k
    }
  }
}

const ThresholdStatement: FunctionComponent<Props> = ({
  status,
  threshold,
  onChangeType,
  onDelete,
  children,
}) => {
  const selected = optionSelector(threshold.type)
  const onChangeThresholdType = (option: string) => {
    onChangeType(dropdownOptions[option])
  }

  return (
    <Panel backgroundColor={InfluxColors.Castle}>
      <DismissButton
        color={ComponentColor.Default}
        type={ButtonType.Button}
        testID={'dismiss--button'}
        onClick={onDelete}
      />

      <PanelBody>
        <FlexBox direction={FlexDirection.Column} margin={ComponentSize.Small}>
          <FlexBox
            direction={FlexDirection.Row}
            margin={ComponentSize.Small}
            stretchToFitWidth={true}
          >
            <TextBlock text={'when value'} />
            <FlexBoxChild grow={2}>
              <SelectDropdown
                selectedOption={selected}
                options={Object.keys(dropdownOptions)}
                onSelect={onChangeThresholdType}
              />
            </FlexBoxChild>
          </FlexBox>

          <FlexBox
            direction={FlexDirection.Row}
            margin={ComponentSize.Small}
            stretchToFitWidth={true}
          >
            {children}
            <TextBlock text={'set status to'} />
            <TextBlock
              text={status}
              backgroundColor={CHECK_STATUS_COLORS[status]}
              testID={'threshold-level-text-block'}
            />
          </FlexBox>
        </FlexBox>
      </PanelBody>
    </Panel>
  )
}

export default ThresholdStatement
