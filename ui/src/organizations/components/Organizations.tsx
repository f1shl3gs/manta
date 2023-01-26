// Libraries
import React, {FunctionComponent, ReactNode, useEffect, useState} from 'react'

// Components
import Nav from 'src/layout/Nav'
import PageSpinner from 'src/shared/components/PageSpinner'

// Types
import {RemoteDataState} from '@influxdata/clockface'

// Hooks
import {useDispatch, useSelector} from 'react-redux'

// Actions
import {getOrgs} from 'src/organizations/actions/thunks'
import {getOrg, getOrgs as selectOrgs} from 'src/organizations/selectors'
import {setOrg} from 'src/organizations/actions'

interface Props {
  children: ReactNode
}

// just get all organizations
const Organizations: FunctionComponent<Props> = ({children}) => {
  const dispatch = useDispatch()
  const [loading, setLoading] = useState(RemoteDataState.Loading)
  const orgs = useSelector(selectOrgs)
  const org = useSelector(getOrg)

  useEffect(() => {
    if (orgs.length === 0) {
      dispatch(getOrgs())
      return
    }

    setLoading(RemoteDataState.Done)
    if (!org) {
      dispatch(setOrg(orgs[0]))
    }
  }, [org, orgs, dispatch])

  return (
    <PageSpinner loading={loading}>
      <Nav />
      <>{children}</>
    </PageSpinner>
  )
}

export default Organizations
