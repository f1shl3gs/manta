// Libraries
import constate from 'constate'
import {useEffect, useState} from 'react'

// Types
import {Organization} from 'src/types/Organization'
import {matchPath, useNavigate} from 'react-router-dom'

const [OrganizationsProvider, useOrganizations, useOrganization] = constate(
  (state: {organizations: Organization[]}) => {
    const organizations = state.organizations
    const [current, setCurrent] = useState(0)
    const navigate = useNavigate()

    useEffect(() => {
      if (matchPath('/orgs/*', window.location.pathname) === null) {
        navigate(`/orgs/${organizations[current].id}`)
      }
    }, [organizations, current, navigate])

    return {
      current,
      organizations,
      setCurrent,
    }
  },
  value => value,
  value => value.organizations[value.current]
)

export {OrganizationsProvider, useOrganizations, useOrganization}
