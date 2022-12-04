// Libraries
import React, {useRef, FunctionComponent} from 'react'
import dayjs from 'dayjs'
import {useDispatch, useSelector} from 'react-redux'

// Components
import {
  Appearance,
  Dropdown,
  Popover,
  PopoverInteraction,
  PopoverPosition,
} from '@influxdata/clockface'
import DateRangePicker from 'src/shared/components/DateRangePicker/DateRangePicker'

// Hooks

// Constants
import {
  CUSTOM_TIME_RANGE_LABEL,
  SELECTABLE_TIME_RANGES,
  TIME_RANGE_FORMAT,
} from 'src/constants/timeRange'

// Types
import {TimeRange} from 'src/types/timeRanges'
import {getTimeRange} from 'src/shared/selectors/timeRange'

// Actions
import {setTimeRange} from 'src/shared/actions/timeRange'

const getTimeRangeLabel = (timeRange: TimeRange): string => {
  if (timeRange.type === 'selectable-duration') {
    return timeRange.label
  }

  if (timeRange.type === 'duration') {
    return timeRange.lower
  }

  if (timeRange.type === 'custom') {
    return `${dayjs(timeRange.lower).format(TIME_RANGE_FORMAT)} - ${dayjs(
      timeRange.upper
    ).format(TIME_RANGE_FORMAT)}`
  }

  return 'unknown'
}

const TimeRangeDropdown: FunctionComponent = () => {
  const timeRange = useSelector(getTimeRange)
  const dispatch = useDispatch()
  const dropdownRef = useRef<HTMLDivElement>(null)
  const timeRangeLabel = getTimeRangeLabel(timeRange)

  const dropdownWidth = (): number => {
    if (timeRange.type === 'custom') {
      return 280
    }

    return 110
  }

  return (
    <>
      <Popover
        appearance={Appearance.Outline}
        position={PopoverPosition.ToTheLeft}
        triggerRef={dropdownRef}
        visible={false}
        showEvent={PopoverInteraction.None}
        hideEvent={PopoverInteraction.None}
        distanceFromTrigger={8}
        enableDefaultStyles={false}
        contents={() => (
          <DateRangePicker
            timeRange={timeRange}
            onSetTimeRange={tr => {
              setTimeRange(tr)
            }}
            onClose={() => {
              console.log('close')
            }}
          />
        )}
      />

      <div ref={dropdownRef}>
        <Dropdown
          style={{width: `${dropdownWidth()}px`}}
          button={(active, onClick) => (
            <Dropdown.Button active={active} onClick={onClick}>
              {timeRangeLabel}
            </Dropdown.Button>
          )}
          menu={onCollapse => (
            <Dropdown.Menu
              onCollapse={onCollapse}
              style={{width: `${dropdownWidth() + 50}px`}}
            >
              <Dropdown.Divider
                key={'Time Range'}
                text={'Time Range'}
                id={'Time Range'}
              />

              <Dropdown.Item
                key={CUSTOM_TIME_RANGE_LABEL}
                value={CUSTOM_TIME_RANGE_LABEL}
                id={CUSTOM_TIME_RANGE_LABEL}
                selected={true}
                onClick={() => console.log('onclick')}
              >
                {CUSTOM_TIME_RANGE_LABEL}
              </Dropdown.Item>

              {SELECTABLE_TIME_RANGES.map(item => {
                const {label} = item
                const handleClick = (): void => {
                  dispatch(setTimeRange(item))
                }

                return (
                  <Dropdown.Item
                    key={label}
                    value={label}
                    id={label}
                    selected={label === timeRangeLabel}
                    onClick={handleClick}
                  >
                    {label}
                  </Dropdown.Item>
                )
              })}
            </Dropdown.Menu>
          )}
        />
      </div>
    </>
  )
}

export default TimeRangeDropdown
