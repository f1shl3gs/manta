import React, {FunctionComponent, useState} from 'react'
import {TreeNav} from '@influxdata/clockface'
import OrganizationsSwitcher from './OrganizationsSwitcher'
import {useUser} from '../shared/components/useAuthentication'
import {useOrganization} from './useOrganizations'
import useFetch from '../shared/useFetch'
import {useNavigate} from 'react-router-dom';

const UserWidget: FunctionComponent = () => {
  const [switcherVisible, setSwitcherVisible] = useState(false)
  const user = useUser()
  const org = useOrganization()
  const navigate = useNavigate()
  const {run: logout} = useFetch(`/api/v1/signout`, {
    method: 'DELETE',
    onSuccess: _ => navigate(`/signin`)
  })

  return (
    <div>
      <OrganizationsSwitcher
        visible={switcherVisible}
        dismiss={() => setSwitcherVisible(false)}
      />

      <TreeNav.User
        username={user.name}
        team={org.name}
        testID={'tree-nav-user'}
      >
        <TreeNav.SubHeading label="Team" />
        <TreeNav.UserItem id="members" label="Members" />
        <TreeNav.UserItem id="about" label="About" />

        <TreeNav.SubHeading label={user.name} lowercase />
        <TreeNav.UserItem
          id="switch"
          label="Switch organization"
          onClick={() => setSwitcherVisible(true)}
        />
        <TreeNav.UserItem
          id="logout"
          label="Logout"
          testID={'user-logout'}
          onClick={() => logout()}
        />
      </TreeNav.User>
    </div>
  )
}

export default UserWidget
