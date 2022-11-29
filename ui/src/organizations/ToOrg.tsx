import React, {FunctionComponent, useEffect} from 'react'
import {useNavigate} from 'react-router-dom'
import {useOrganization} from 'src/organizations/useOrganizations'

const ToOrg: FunctionComponent = () => {
  const navigate = useNavigate()
  const {id} = useOrganization()

  useEffect(() => {
    navigate(`/orgs/${id}`)
  }, [id, navigate])

  return <></>
}

export default ToOrg
