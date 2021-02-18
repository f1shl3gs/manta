import {useState} from 'react'
import constate from 'constate'

import {Organization} from 'types/Organization'

type State = {
  initialOrg: Organization
}

const [OrgProvider, useOrgID] = constate(
  (state: State) => {
    const [org] = useState<Organization>(state.initialOrg)
    return {org}
  },
  value => value.org.id
)

export {OrgProvider, useOrgID}
