// Libraries
import React, {useEffect, useState} from 'react'

// Components
import {
  ClickOutside,
  ComponentStatus,
  DropdownDivider,
  DropdownItem,
  DropdownMenu,
  Input,
} from '@influxdata/clockface'

// Constants
const SUGGESTION_CLASS = 'duration-input--suggestion'

// Utils
export const isDurationParsable = (duration: string): boolean => {
  const durationRegExp = /^(([0-9]+)(y|mo|w|d|h|ms|s|m|us|µs|ns))+$/g

  // warning! Using string.match(regex) here instead of regex.test(string) because regex.test() modifies the regex object, and can lead to unexpected behavior

  return !!duration.match(durationRegExp)
}

interface Props {
  suggestions: string[]
  onSubmit: (v: string) => void
  value: string
  placeholder?: string
  submitInvalid?: boolean
  showDivider?: boolean
  testID?: string
  validFunction?: (v: string) => boolean
  status?: ComponentStatus
}

const DurationInput: React.FC<Props> = props => {
  const {
    suggestions,
    onSubmit,
    value,
    placeholder,
    status: controlledStatus,
    submitInvalid = true,
    showDivider = true,
    testID = 'duration-input',
    validFunction = _ => false,
  } = props

  const [focused, setFocused] = useState(false)
  const [inputValue, setInputValue] = useState(value)

  useEffect(() => {
    setInputValue(value)
  }, [value, setInputValue])

  const handleClickSuggestion = (suggestion: string) => {
    setInputValue(suggestion)
    onSubmit(suggestion)
    setFocused(false)
  }

  // @ts-ignore
  const handleClickOutside = e => {
    const didClickSuggestion =
      e.target.classList.contains(SUGGESTION_CLASS) ||
      e.target.parentNode.classList.contains(SUGGESTION_CLASS)

    if (!didClickSuggestion) {
      setFocused(false)
    }
  }

  const isValid = (v: string): boolean =>
    isDurationParsable(v) || validFunction(v)

  let inputStatus = controlledStatus || ComponentStatus.Default

  if (inputStatus === ComponentStatus.Default && !isValid(inputValue)) {
    inputStatus = ComponentStatus.Error
  }

  const onChange = (v: string) => {
    setInputValue(v)

    if (submitInvalid || (!submitInvalid && isValid(v))) {
      onSubmit(v)
    }
  }

  return (
    <div className={'duration-input'}>
      <ClickOutside onClickOutside={handleClickOutside}>
        <Input
          placeholder={placeholder}
          value={inputValue}
          status={inputStatus}
          onChange={e => onChange(e.target.value)}
          onFocus={() => setFocused(true)}
          testID={testID}
        />
      </ClickOutside>
      {focused && (
        <DropdownMenu className={'duration-input--menu'} noScrollX={true}>
          {showDivider && <DropdownDivider text={'Examples'} />}
          {suggestions.map(s => (
            <DropdownItem
              key={s}
              value={s}
              className={SUGGESTION_CLASS}
              onClick={handleClickSuggestion}
            >
              {s}
            </DropdownItem>
          ))}
        </DropdownMenu>
      )}
    </div>
  )
}

export default DurationInput
