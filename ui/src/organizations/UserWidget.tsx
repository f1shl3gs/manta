import React, {FunctionComponent, useState} from 'react'
import {TreeNav} from '@influxdata/clockface'
import OrganizationsSwitcher from './OrganizationsSwitcher'
import {useUser} from '../shared/components/useAuthentication'
import {useOrganization} from './useOrganizations'

const UserWidget: FunctionComponent = () => {
  const [switcherVisible, setSwitcherVisible] = useState(false)
  const user = useUser()
  const org = useOrganization()

  return (
    <div>
      <OrganizationsSwitcher
        visible={switcherVisible}
        dismiss={() => setSwitcherVisible(false)}
      />

      <TreeNav.User username={user.name} team={org.name} testID={'tree-nav-user'}>
        <TreeNav.SubHeading label="Team" />
        <TreeNav.UserItem id="members" label="Members" />
        <TreeNav.UserItem id="about" label="About" />

        <TreeNav.SubHeading label={user.name} lowercase />
        <TreeNav.UserItem
          id="switch"
          label="Switch organization"
          onClick={() => setSwitcherVisible(true)}
        />
        <TreeNav.UserItem id="logout" label="Logout" />
      </TreeNav.User>
    </div>
  )
}

export default UserWidget
