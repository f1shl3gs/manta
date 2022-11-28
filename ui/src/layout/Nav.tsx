// Libraries
import React, {FunctionComponent} from 'react'
import {connect, ConnectedProps} from 'react-redux'

// Components
import {
  Icon,
  IconFont,
  InfluxDBCloudLogo,
  TreeNav,
  TreeNavItem,
  TreeNavSubItem,
  TreeNavSubMenu,
} from '@influxdata/clockface'
import {Link, useNavigate} from 'react-router-dom'
import UserWidget from 'src/organizations/UserWidget'

// Hooks
import {useOrganization} from 'src/organizations/useOrganizations'

// Actions
import { toggleNavBarState } from 'src/shared/actions/app'

// Types
import {AppState} from 'src/types/stores'

const getNavItemActivation = (
  keywords: string[],
  location: string
): boolean => {
  const ignoreOrgAndOrgId = 3
  const parentPath = location.split('/').slice(ignoreOrgAndOrgId)

  if (!parentPath.length) {
    parentPath.push('me')
  }

  return keywords.some(path => parentPath.includes(path))
}

interface NavItemLink {
  type: 'link' | 'href'
  location: string
}

interface NavItem {
  id: string
  testID: string
  label: string
  shortLabel?: string
  icon: IconFont
  activeKeywords: string[]
  link: NavItemLink
  menu?: NavSubItem[]
}

interface NavSubItem {
  id: string
  testID: string
  label: string
  link: NavItemLink
}

const generateNavItems = (orgId: string): NavItem[] => {
  const orgPrefix = `/orgs/${orgId}`

  return [
    {
      id: 'data',
      testID: 'nav-item-data',
      label: 'Data',
      icon: IconFont.Upload_New,
      shortLabel: 'Data',
      activeKeywords: ['vertex', 'data'],
      link: {
        type: 'link',
        location: `${orgPrefix}/data`,
      },
      menu: [
        {
          id: 'vertex',
          testID: 'nav-item-data--vertex',
          label: 'Vertex',
          link: {
            type: 'link',
            location: `${orgPrefix}/data/vertex`,
          },
        },
        {
          id: 'config',
          testID: 'nav-item-data--config',
          label: 'Config',
          link: {
            type: 'link',
            location: `${orgPrefix}/data/config`,
          },
        },
        {
          id: 'scrape',
          testID: 'nav-item-data--scrape',
          label: 'Scrape',
          link: {
            type: 'link',
            location: `${orgPrefix}/data/scrape`,
          },
        },
      ],
    },
    {
      id: 'explore',
      testID: 'nav-item-explore',
      label: 'Explore',
      icon: IconFont.Search_New,
      shortLabel: 'Explore',
      activeKeywords: ['explore'],
      link: {
        type: 'link',
        location: `${orgPrefix}/explore`,
      },
    },
    {
      id: 'dashboards',
      testID: 'nav-item-dashboard',
      label: 'Dashboards',
      icon: IconFont.GraphLine_New,
      shortLabel: 'Dashboards',
      activeKeywords: ['dashboards'],
      link: {
        type: 'link',
        location: `${orgPrefix}/dashboards`,
      },
    },
    {
      id: 'todo',
      testID: 'todo',
      label: 'Todo',
      icon: IconFont.Annotate_New,
      shortLabel: 'Todo',
      activeKeywords: ['todo'],
      link: {
        type: 'link',
        location: `${orgPrefix}/todo`,
      },
    },
    {
      id: 'settings',
      testID: 'nav-item-settings',
      label: 'Settings',
      icon: IconFont.CogOutline_New,
      activeKeywords: ['settings', 'members', 'secrets'],
      link: {
        type: 'link',
        location: `${orgPrefix}/settings`,
      },
      menu: [
        {
          id: 'members',
          testID: 'members',
          label: 'Members',
          link: {
            type: 'link',
            location: `${orgPrefix}/settings/members`,
          },
        },
        {
          id: 'secrets',
          testID: 'secrets',
          label: 'Secrets',
          link: {
            type: 'link',
            location: `${orgPrefix}/settings/secrets`,
          },
        },
      ],
    },
  ]
}

type Props = ConnectedProps<typeof connector>

const Nav: FunctionComponent<Props> = ({navbarState, toggleNavBarState}) => {
  const {id: orgId} = useOrganization()
  const navigate = useNavigate()
  const navItems = generateNavItems(orgId)

  return (
    <TreeNav
      expanded={navbarState}
      onToggleClick={toggleNavBarState}
      headerElement={
        <TreeNav.Header
          id="home"
          label={<InfluxDBCloudLogo cloud={true} />}
          onClick={() => navigate(`/orgs/${orgId}`)}
          icon={<Icon glyph={IconFont.CuboSolid} />}
        />
      }
      userElement={<UserWidget />}
    >
      {navItems.map(item => {
        const linkElement = (classname: string): JSX.Element => {
          if (item.link.type === 'href') {
            return <a href={item.link.location} className={classname} />
          }

          return <Link to={item.link.location} className={classname} />
        }

        return (
          <TreeNavItem
            key={item.id}
            id={item.id}
            testID={item.testID}
            icon={<Icon glyph={item.icon} />}
            label={item.label}
            shortLabel={item.shortLabel || ''}
            active={getNavItemActivation(
              item.activeKeywords,
              window.location.pathname
            )}
            linkElement={linkElement}
          >
            {Boolean(item.menu) && (
              <TreeNavSubMenu>
                {item.menu?.map(menuItem => {
                  const linkElement = (classname: string): JSX.Element => {
                    if (menuItem.link.type === 'href') {
                      return (
                        <a
                          href={menuItem.link.location}
                          className={classname}
                        />
                      )
                    }

                    return (
                      <Link to={menuItem.link.location} className={classname} />
                    )
                  }

                  return (
                    <TreeNavSubItem
                      id={menuItem.id}
                      key={menuItem.id}
                      testID={menuItem.testID}
                      label={menuItem.label}
                      linkElement={linkElement}
                    />
                  )
                })}
              </TreeNavSubMenu>
            )}
          </TreeNavItem>
        )
      })}
    </TreeNav>
  )
}

const mstp = ({
  app: {
    persisted: {navbarState}
  }
}: AppState) => {
  return {
    navbarState: navbarState === 'expanded'
  }
}

const mdtp = {
  toggleNavBarState
}

const connector = connect(mstp, mdtp)

export default connector(Nav)
