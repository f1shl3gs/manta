// Libraries
import React, {useRef, useState} from 'react'
import moment from 'moment'

// Components
import {
  Appearance,
  Dropdown,
  Popover,
  PopoverInteraction,
  PopoverPosition,
} from '@influxdata/clockface'
import DateRangePicker from 'shared/components/DateRangePicker/DateRangePicker'

// Hooks
import {useTimeRange} from 'shared/useTimeRange'

// Types
import {TimeRange} from 'types/TimeRanges'

// Constants
import {
  CUSTOM_TIME_RANGE_LABEL,
  SELECTABLE_TIME_RANGES,
  TIME_RANGE_FORMAT,
} from 'constants/timeRange'

const getTimeRangeLabel = (timeRange: TimeRange): string => {
  if (timeRange.type === 'selectable-duration') {
    return timeRange.label
  }

  if (timeRange.type === 'duration') {
    return timeRange.lower
  }

  if (timeRange.type === 'custom') {
    return `${moment(timeRange.lower).format(TIME_RANGE_FORMAT)} - ${moment(
      timeRange.upper
    ).format(TIME_RANGE_FORMAT)}`
  }

  return 'unknown'
}

const TimeRangeDropdown: React.FC = () => {
  const {timeRange, setTimeRange} = useTimeRange()
  const dropdownRef = useRef<HTMLDivElement>(null)
  const [visible, setVisible] = useState(false)
  const timeRangeLabel = getTimeRangeLabel(timeRange)

  const dropdownWidth = (): number => {
    if (timeRange.type === 'custom') {
      return 250
    }

    return 100
  }

  return (
    <>
      <Popover
        appearance={Appearance.Outline}
        position={PopoverPosition.ToTheLeft}
        triggerRef={dropdownRef}
        visible={visible}
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

                return (
                  <Dropdown.Item
                    key={label}
                    value={label}
                    id={label}
                    selected={label === timeRangeLabel}
                    onClick={() => {
                      setTimeRange(item)
                    }}
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
