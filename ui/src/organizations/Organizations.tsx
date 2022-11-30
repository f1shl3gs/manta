// Libraries
import React, {FunctionComponent, ReactNode, useEffect, useState} from 'react'

// Components
import Nav from 'src/layout/Nav'
import PageSpinner from 'src/shared/components/PageSpinner'

// Hooks
import {useDispatch, useSelector} from 'react-redux'

// Actions
import {getOrgs} from 'src/organizations/actions/thunks'
import {getOrgs as selectOrgs} from 'src/organizations/selectors'
import {RemoteDataState} from '@influxdata/clockface'
import {setOrg} from 'src/organizations/actions'

interface Props {
  children: ReactNode
}

// just get all organizations
const Organizations: FunctionComponent<Props> = ({children}) => {
  const dispatch = useDispatch()
  const [loading, setLoading] = useState(RemoteDataState.Loading)
  const orgs = useSelector(selectOrgs)

  useEffect(() => {
    if (orgs.length === 0) {
      dispatch(getOrgs())
      return
    }

    setLoading(RemoteDataState.Done)
    dispatch(setOrg(orgs[0]))
  }, [orgs, dispatch])

  return (
    <PageSpinner loading={loading}>
      <Nav />
      <>{children}</>
    </PageSpinner>
  )
}

export default Organizations
