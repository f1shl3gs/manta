// Libraries
import constate from 'constate'
import {useCallback, useEffect, useState} from 'react'
import {useParams} from 'react-router-dom'

// Hooks
import {useOrgID} from 'shared/useOrg'

// types
import {Check, CheckStatusLevel, Condition} from '../../types/Check'
import {useFetch} from 'shared/useFetch'

// Utils
import {RemoteDataState} from '@influxdata/clockface'
import {
  defaultErrorNotification,
  useNotification,
} from 'shared/notification/useNotification'

interface CheckUpdate {
  name: string
}

// todo: fix the status properties
const defaultCheck: Check = {
  lastRunStatus: '',
  latestCompleted: '',
  latestFailure: '',
  latestScheduled: '',
  latestSuccess: '',
  conditions: [],
  created: '',
  id: '',
  orgID: '',
  status: '',
  updated: '',
  name: '',
  desc: '',
  expr: '',
  labels: [],
  cron: '@every 1m',
}

const [CheckProvider, useCheck] = constate(
  () => {
    const {id} = useParams<{id: string}>()
    const [tab, setTab] = useState('query')
    const [check, setCheck] = useState<Check>(defaultCheck)
    const {notify} = useNotification()
    const orgID = useOrgID()
    const [remoteDataState, setRemoteDataState] = useState(
      RemoteDataState.NotStarted
    )

    const {get, post} = useFetch(`/api/v1/checks/${id}`, {})
    const {put} = useFetch(`/api/v1/checks?orgID=${orgID}`, {})

    useEffect(() => {
      if (id === 'new') {
        setRemoteDataState(RemoteDataState.Done)
        return
      }

      setRemoteDataState(RemoteDataState.Loading)
      get()
        .then(data => {
          // todo: dummy
          if (data.labels === null) {
            data.labels = []
          }

          setCheck(data)
          setRemoteDataState(RemoteDataState.Done)
        })
        .catch(err => {
          setRemoteDataState(RemoteDataState.Error)
          notify({
            ...defaultErrorNotification,
            message: `Fetch Check failed, err: ${err.message}`,
          })
        })
    }, [get, id, notify])

    const updateCheck = (udp: CheckUpdate) => {
      // @ts-ignore
      setCheck(prev => {
        return {
          ...prev,
          ...udp,
        }
      })
    }

    const onSave = useCallback(() => {
      if (id === 'new') {
        return put(check)
      } else {
        return post(check)
      }
    }, [check, id, post, put])

    const onRename = useCallback((name: string) => {
      setCheck(prev => {
        return {
          ...prev,
          name,
        }
      })
    }, [])

    const onExprUpdate = useCallback(
      (expr: string) => {
        // promql editor will call this when mount
        // it must be stopped
        if (check.expr === expr) {
          return
        }

        setCheck(prev => {
          return {
            ...prev,
            expr,
          }
        })
      },
      [check]
    )

    const onAddCondition = useCallback((level: CheckStatusLevel) => {
      setCheck(prevState => {
        const {conditions} = prevState

        conditions.push({
          status: level,
          threshold: {
            type: 'gt',
            value: 0,
          },
        })
        return {
          ...prevState,
          conditions,
        }
      })
    }, [])

    const onSetCron = useCallback((cron: string) => {
      setCheck(prev => {
        return {
          ...prev,
          cron,
        }
      })
    }, [])

    const onSetOffset = useCallback((offset: string) => {
      setCheck(prev => {
        return {
          ...prev,
          offset,
        }
      })
    }, [])

    const onChangeCondition = useCallback((condition: Condition) => {
      setCheck(prev => {
        const {conditions} = prev

        const next = conditions.map(item => {
          if (item.status !== (condition.status as CheckStatusLevel)) {
            return item
          }

          return condition
        })

        return {
          ...prev,
          conditions: next,
        }
      })
    }, [])

    return {
      ...check,
      tab,
      setTab,
      onSave,
      onRename,
      onSetCron,
      onSetOffset,
      onExprUpdate,
      onAddCondition,
      onChangeCondition,
      remoteDataState,
      updateCheck,
    }
  },
  value => value
)

export {CheckProvider, useCheck}
