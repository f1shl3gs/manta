// Libraries
import React, {FunctionComponent, useCallback, useState} from 'react'
import {useSearchParams} from 'react-router-dom'

// Componetns
import {
  ComponentStatus,
  Dropdown,
  DropdownButton,
  DropdownDivider,
  DropdownItem,
  DropdownMenu,
  IconFont,
} from '@influxdata/clockface'

// Types
import {
  AutoRefresh,
  AutoRefreshOption,
  AutoRefreshOptionType,
  AutoRefreshStatus,
} from 'src/types/autoRefresh'

const autoRefreshOptions: AutoRefreshOption[] = [
  {
    id: 'refresh',
    type: AutoRefreshOptionType.Header,
    label: 'Refresh',
    seconds: 0,
  },
  {
    id: 'pause',
    type: AutoRefreshOptionType.Option,
    label: 'Pause',
    seconds: 0,
  },
  {
    id: '10s',
    type: AutoRefreshOptionType.Option,
    label: '5s',
    seconds: 5,
  },
  {
    id: '15s',
    type: AutoRefreshOptionType.Option,
    label: '15s',
    seconds: 15,
  },
  {
    id: '30s',
    type: AutoRefreshOptionType.Option,
    label: '30s',
    seconds: 30,
  },
  {
    id: '60s',
    type: AutoRefreshOptionType.Option,
    label: '60s',
    seconds: 60,
  },
  {
    id: '5m',
    type: AutoRefreshOptionType.Option,
    label: '5m',
    seconds: 5 * 60,
  },
  {
    id: '15m',
    type: AutoRefreshOptionType.Option,
    label: '15m',
    seconds: 15 * 60,
  },
]

const dropdownIcon = (autoRefresh: AutoRefresh): IconFont => {
  if (autoRefresh.status === AutoRefreshStatus.Paused) {
    return IconFont.Pause
  }

  return IconFont.Refresh_New
}

const dropdownStatus = (autoRefresh: AutoRefresh): ComponentStatus => {
  if (autoRefresh.status === AutoRefreshStatus.Disabled) {
    return ComponentStatus.Disabled
  }

  return ComponentStatus.Default
}

const AutoRefreshDropdown: FunctionComponent = () => {
  const [_, setParams] = useSearchParams()
  const {autoRefresh, setAutoRefresh} = useAutoRefresh()
  const [selected, setSelected] = useState(() => {
    const opt = autoRefreshOptions.find(
      opt => opt.seconds === autoRefresh.interval
    )
    if (opt === undefined) {
      return autoRefreshOptions[3]
    }

    return opt
  })

  const onSelectAutoRefreshOption = useCallback(
    (opt: AutoRefreshOption) => {
      setSelected(opt)
      setAutoRefresh({
        status:
          opt.seconds !== 0
            ? AutoRefreshStatus.Active
            : AutoRefreshStatus.Paused,
        interval: opt.seconds,
      })
      setParams(prev => {
        prev.set('interval', `${opt.seconds}s`)
        return prev
      })
    },
    [setAutoRefresh, setParams]
  )

  return (
    <>
      <Dropdown
        button={(active, onClick) => (
          <DropdownButton
            active={active}
            onClick={onClick}
            status={dropdownStatus(autoRefresh)}
            icon={dropdownIcon(autoRefresh)}
          >
            {selected.label}
          </DropdownButton>
        )}
        menu={onCollapse => (
          <DropdownMenu onCollapse={onCollapse}>
            {autoRefreshOptions.map(option => {
              if (option.type === AutoRefreshOptionType.Header) {
                return (
                  <DropdownDivider
                    key={option.id}
                    id={option.id}
                    text={option.label}
                  />
                )
              }

              return (
                <DropdownItem
                  key={option.id}
                  id={option.id}
                  value={option}
                  onClick={onSelectAutoRefreshOption}
                >
                  {option.label}
                </DropdownItem>
              )
            })}
          </DropdownMenu>
        )}
      />
    </>
  )
}

export default AutoRefreshDropdown
