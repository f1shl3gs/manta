// Libraries
import React, {FunctionComponent, useState} from 'react'
import {useNavigate} from 'react-router-dom'
import {useDispatch, useSelector} from 'react-redux'

// Components
import {TreeNav} from '@influxdata/clockface'
import OrganizationsSwitcher from 'src/organizations/components/OrganizationsSwitcher'

// Actions
import { signout } from 'src/me/actions/thunks'

// Selectors
import {useOrg} from 'src/organizations/selectors'
import {getMeName} from 'src/me/selectors'

const UserWidget: FunctionComponent = () => {
  const [switcherVisible, setSwitcherVisible] = useState(false)
  const username = useSelector(getMeName)
  const org = useOrg()
  const dispatch = useDispatch()
  const navigate = useNavigate()

  const handleSignout = () => {
    dispatch(signout())
  }

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
        username={username}
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
          onClick={handleSignout}
        />
      </TreeNav.User>
    </div>
  )
}

export default React.memo(UserWidget)
