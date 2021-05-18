import constate from 'constate'
import {useEffect, useState} from 'react'
import {Organization} from '../types/Organization'
import {RemoteDataState} from '@influxdata/clockface'

const [OrgsProvider, useOrgs] = constate(
  () => {
    const [orgs, setOrgs] = useState<Organization[]>([])
    const [loading, setLoading] = useState(RemoteDataState.NotStarted)

    useEffect(() => {
      setLoading(RemoteDataState.Loading)
      fetch(`/api/v1/orgs`)
        .then(resp => resp.json())
        .then(data => {
          setOrgs(data)
          setLoading(RemoteDataState.Done)
        })
        .catch(err => {
          setLoading(RemoteDataState.Error)
        })
    }, [])

    return {
      orgs,
      loading,
    }
  },
  values => values
)

export {OrgsProvider, useOrgs}
