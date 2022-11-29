import React, {FunctionComponent, useState} from 'react'
import {TreeNav} from '@influxdata/clockface'
import OrganizationsSwitcher from 'src/organizations/OrganizationsSwitcher'
import {useUser} from 'src/shared/components/useAuthentication'
import useFetch from 'src/shared/useFetch'
import {useNavigate} from 'react-router-dom'
import {useOrg} from 'src/organizations/selectors'

const UserWidget: FunctionComponent = () => {
  const [switcherVisible, setSwitcherVisible] = useState(false)
  const user = useUser()
  const org = useOrg()
  const navigate = useNavigate()
  const {run: logout} = useFetch(`/api/v1/signout`, {
    method: 'DELETE',
    onSuccess: _ => navigate(`/signin`),
  })

  return (
    <div>
      {/*
        TODO: this is a dummy operation, but it did reduce re-render
      */}
      {switcherVisible && (
        <OrganizationsSwitcher
          visible={switcherVisible}
          dismiss={() => setSwitcherVisible(false)}
        />
      )}

      <TreeNav.User
        username={user.name}
        team={org.name}
        testID={'tree-nav-user'}
      >
        <TreeNav.UserItem id="members" label="Members" />
        <TreeNav.UserItem id="about" label="About" />

        <TreeNav.UserItem
          id="switch"
          label="Switch organization"
          testID="switch organization"
          onClick={() => setSwitcherVisible(true)}
        />
        <TreeNav.UserItem
          id={'create-org'}
          label="Create organization"
          testID="create-org"
          onClick={() => navigate('/orgs/new')}
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

export default React.memo(UserWidget)
