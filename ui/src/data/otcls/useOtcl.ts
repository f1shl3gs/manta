import {useCallback, useEffect, useState} from 'react'
import constate from 'constate'
import {RemoteDataState} from '@influxdata/clockface'
import {useHistory, useParams} from 'react-router-dom'
import {
  defaultErrorNotification,
  defaultSuccessNotification,
  useNotification,
} from 'shared/notification/useNotification'
import {useOrgID} from '../../shared/useOrg'

const [OtclProvider, useOtcl] = constate(
  () => {
    const [loading, setLoading] = useState(RemoteDataState.NotStarted)
    const [otcl, setOtcl] = useState({
      id: '',
      name: '',
      desc: '',
      content: '',
    })
    const {notify} = useNotification()
    const {id} = useParams<{id: string}>()
    const history = useHistory()
    const orgID = useOrgID()

    useEffect(() => {
      if (id === 'new') {
        setLoading(RemoteDataState.Done)
        return
      }

      setLoading(RemoteDataState.Loading)
      fetch(`/api/v1/otcls/${id}`, {
        headers: {
          Accept: 'application/json',
        },
      })
        .then(resp => resp.json())
        .then(data => {
          setOtcl(data)
          setLoading(RemoteDataState.Done)
        })
        .catch(err => {
          notify({
            ...defaultErrorNotification,
            message: `Fetch otcl failed, err: ${err.message}`,
          })
          setLoading(RemoteDataState.Error)
        })

      // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    const onSave = useCallback(() => {
      // create
      if (id === 'new') {
        fetch(`/api/v1/otcls`, {
          method: 'POST',
          body: JSON.stringify({
            orgID,
            name: otcl.name,
            content: otcl.content,
          }),
        })
          .then(() => {
            notify({
              ...defaultSuccessNotification,
              message: `Create new Otcl success`,
            })
            history.push(`/orgs/${orgID}/data/otcls`)
          })
          .catch(err => {
            notify({
              ...defaultErrorNotification,
              message: `Create new Otcl failed, err: ${err.message}`,
            })
          })
        return
      }

      // update
      fetch(`/api/v1/otcls/${id}`, {
        method: 'PATCH',
        body: JSON.stringify({
          name: otcl.name,
          content: otcl.content,
        }),
      })
        .then(() => {
          notify({
            ...defaultSuccessNotification,
            message: `Update otcl "${otcl.name}" success`,
          })
          history.push(`/orgs/${orgID}/data/otcls`)
        })
        .catch(err => {
          notify({
            ...defaultErrorNotification,
            message: `Update otcl "${otcl.name}" failed, err: ${err.message}`,
          })
        })
    }, [history, id, notify, orgID, otcl])

    const onRename = useCallback(
      (name: string) => {
        setOtcl({
          ...otcl,
          name,
        })
      },
      [otcl, setOtcl]
    )

    const onContentChange = useCallback(
      content => {
        setOtcl({
          ...otcl,
          content,
        })
      },
      [otcl]
    )

    return {
      otcl,
      loading,
      onSave,
      onRename,
      onContentChange,
    }
  },
  value => value
)

export {OtclProvider, useOtcl}
