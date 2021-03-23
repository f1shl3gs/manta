// Libraries
import React from 'react'

// Components
import {
  ButtonType,
  ComponentColor,
  ComponentSize,
  DismissButton,
  FlexBox,
  FlexDirection,
  InfluxColors,
  Panel,
  PanelBody,
  SelectDropdown,
  TextBlock,
} from '@influxdata/clockface'

// Types
import {CheckStatusLevel, Threshold, ThresholdType} from 'types/Check'
import {LEVEL_COLORS} from 'constants/level'

interface Props {
  level: CheckStatusLevel
  threshold: Threshold
  changeThresholdType: (to: ThresholdType) => void
}

const dropdownOptions: {[key: string]: ThresholdType} = {
  'is above': 'gt',
  'is below': 'lt',
  'is inside range': 'inside',
  'is outside range': 'outside',
}

const optionSelector = (threshold: Threshold) => {
  const {type} = threshold

  for (const k in dropdownOptions) {
    if (dropdownOptions[k] === type) {
      return dropdownOptions[k]
    }
  }

  return 'unknown'
}

const ThresholdStatement: React.FC<Props> = props => {
  const {children, threshold, level, changeThresholdType} = props
  const selectedOption = optionSelector(threshold)
  const onChangeThresholdType = (option: string) => {
    changeThresholdType(dropdownOptions[option])
  }

  return (
    <Panel backgroundColor={InfluxColors.Castle} testID={'panel'}>
      <DismissButton
        color={ComponentColor.Default}
        onClick={e => console.log(e)}
        testID={'dismiss-button'}
        type={ButtonType.Button}
      />
      <PanelBody testID={'panel--body'}>
        <FlexBox
          direction={FlexDirection.Column}
          margin={ComponentSize.Small}
          testID={'component-spacer'}
        >
          <FlexBox
            direction={FlexDirection.Row}
            margin={ComponentSize.Small}
            stretchToFitWidth
            testID={'component-spacer'}
          >
            <TextBlock testID={'when-value-text-block'} text={'When value '} />
            <FlexBox.Child grow={2} testID={'component-spacer--flex-child'}>
              <SelectDropdown
                options={Object.keys(dropdownOptions)}
                selectedOption={selectedOption}
                onSelect={onChangeThresholdType}
                testID={'select-option-dropdown'}
              />
            </FlexBox.Child>
          </FlexBox>

          <FlexBox
            direction={FlexDirection.Row}
            margin={ComponentSize.Small}
            stretchToFitWidth
            testID={'component-spacer'}
          >
            {children}
            <TextBlock
              testID={'set-status-to-text-block'}
              text={'set status to'}
            />
            <TextBlock
              backgroundColor={LEVEL_COLORS[level]}
              testID={'threshold-level-text-block'}
              text={level}
            />
          </FlexBox>
        </FlexBox>
      </PanelBody>
    </Panel>
  )
}

export default ThresholdStatement
