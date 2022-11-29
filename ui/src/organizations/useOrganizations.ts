// Libraries
import constate from 'constate'
import {useLayoutEffect, useMemo} from 'react'

// Types
import {Organization} from 'src/types/Organization'
import useLocalStorage from 'src/shared/useLocalStorage'

interface State {
  organizations: Organization[]
  refetch: () => void
}

const [OrganizationsProvider, useOrganizations, useOrganization] = constate(
  (state: State) => {
    const {organizations, refetch} = state
    const [current, setCurrent] = useLocalStorage(
      'org',
      organizations[organizations.length - 1]
    )

    // stored organization might outdated, so we must re-store the updated one
    useLayoutEffect(() => {
      const found = organizations.indexOf(current)

      if (found === -1) {
        // not found
        setCurrent(organizations[organizations.length - 1])
      }
    }, [current, setCurrent, organizations])

    return {
      current,
      organizations,
      setCurrent,
      refetch,
    }
  },
  value => useMemo(() => value, [value]),
  value => useMemo(() => value.current, [value.current])
)

export {OrganizationsProvider, useOrganizations, useOrganization}
