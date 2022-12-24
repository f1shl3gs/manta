// Libraries
import React, {FunctionComponent, useCallback, useEffect, useState} from 'react'
import {useSearchParams} from 'react-router-dom'
import {connect, ConnectedProps} from 'react-redux'

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
import {AppState} from 'src/types/stores'

// Actions
import {poll, setAutoRefreshInterval} from 'src/shared/actions/autoRefresh'

// Constants
import {AutoRefreshDropdownOptions} from 'src/shared/constants/autoRefresh'

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

const mstp = (state: AppState) => {
  const {autoRefresh} = state.autoRefresh

  return {
    autoRefresh,
  }
}

const mdtp = {
  setAutoRefreshInterval,
  updateAutoRefresh: poll,
}

const connector = connect(mstp, mdtp)

type Props = ConnectedProps<typeof connector>

const AutoRefreshDropdown: FunctionComponent<Props> = ({
  autoRefresh,
  setAutoRefreshInterval,
  updateAutoRefresh,
}) => {
  const [_, setParams] = useSearchParams()
  const [selected, setSelected] = useState(() => {
    const opt = AutoRefreshDropdownOptions.find(
      opt => opt.seconds === autoRefresh.interval
    )
    if (opt === undefined) {
      return AutoRefreshDropdownOptions[3]
    }

    return opt
  })

  const onSelectAutoRefreshOption = useCallback(
    (opt: AutoRefreshOption) => {
      setSelected(opt)
      setAutoRefreshInterval(opt.seconds)
      setParams(
        prev => {
          prev.set('interval', `${opt.seconds}s`)
          return prev
        },
        {replace: true}
      )
    },
    [setAutoRefreshInterval, setParams]
  )

  useEffect(() => {
    if (autoRefresh.interval === 0) {
      return
    }

    const timer = setInterval(() => {
      if (document.hidden) {
        // tab is not focused, no need to refresh
        return
      }

      updateAutoRefresh()
    }, autoRefresh.interval * 1000)

    return () => {
      clearInterval(timer)
    }
  }, [autoRefresh, updateAutoRefresh])

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
            {AutoRefreshDropdownOptions.map(option => {
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

export default connector(AutoRefreshDropdown)
