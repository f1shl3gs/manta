// Libraries
import React, {FunctionComponent, useCallback, useEffect, useState} from 'react'
import {useSearchParams} from 'react-router-dom'
import {connect, ConnectedProps, useDispatch} from 'react-redux'

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

const mstp = (state: AppState) => {
  const {autoRefresh} = state.autoRefresh

  return {
    autoRefresh,
  }
}

const mdtp = {
  setAutoRefresh: setAutoRefreshInterval,
}

const connector = connect(mstp, mdtp)

type Props = ConnectedProps<typeof connector>

const AutoRefreshDropdown: FunctionComponent<Props> = ({
  autoRefresh,
  setAutoRefresh,
}) => {
  const dispatch = useDispatch()
  const [_, setParams] = useSearchParams()
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
      setAutoRefresh(opt.seconds)
      setParams(
        prev => {
          prev.set('interval', `${opt.seconds}s`)
          return prev
        },
        {replace: true}
      )
    },
    [setAutoRefresh, setParams]
  )

  useEffect(() => {
    dispatch(poll())

    if (autoRefresh.status !== AutoRefreshStatus.Active) {
      return
    }

    const timer = setInterval(() => {
      if (document.hidden) {
        // tab is not focused, no need to refresh
        return
      }

      dispatch(poll())
    }, autoRefresh.interval * 1000)

    return () => {
      clearInterval(timer)
    }
  }, [autoRefresh, dispatch])

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

export default connector(AutoRefreshDropdown)
