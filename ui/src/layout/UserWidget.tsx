import React from 'react'
import {TreeNav} from '@influxdata/clockface'
import {Link, useLocation} from 'react-router-dom'
import {useOrgID} from '../shared/useOrg'

const getNavItemActivation = (
  keywords: string[],
  location: string
): boolean => {
  const ignoreOrgAndOrgID = 3
  const parentPath = location.split('/').slice(ignoreOrgAndOrgID)
  if (!parentPath.length) {
    parentPath.push('/me')
  }

  return keywords.some((path) => parentPath.includes(path))
}

const UserWidget: React.FC = () => {
  // todo: me
  const location = useLocation()
  const orgID = useOrgID()
  const orgPrefix = `/orgs/${orgID}`

  return (
    <TreeNav.User username={'username'} team={'org name'}>
      <TreeNav.UserItem
        id={'members'}
        label={'Members'}
        active={getNavItemActivation(['members'], location.pathname)}
        linkElement={(className) => (
          <Link className={className} to={`${orgPrefix}/about`} />
        )}
      />
      <TreeNav.UserItem
        id={'switch-orgs'}
        label={'Switch Organizations'}
        onClick={() => console.log('switch')}
      />
      <TreeNav.UserItem
        id={'logout'}
        label={'Logout'}
        linkElement={(className) => (
          <Link className={className} to={'/logout'} />
        )}
      />
    </TreeNav.User>
  )
}

export default UserWidget
