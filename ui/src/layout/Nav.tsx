// libraries
import React from 'react'
import {Link} from 'react-router-dom'

// components
import {
  ComponentColor,
  Icon,
  IconFont,
  InfluxDBCloudLogo,
  TreeNav,
} from '@influxdata/clockface'
import {usePresentationMode} from '../shared/usePresentationMode'
import UserWidget from './UserWidget'
import {useOrgID} from '../shared/useOrg'

export interface NavItemLink {
  type: 'link' | 'href'
  location: string
}

const getNavItemActivation = (
  keywords: string[],
  location: string
): boolean => {
  const ignoreOrgAndOrgID = 3
  const parentPath = location.split('/').slice(ignoreOrgAndOrgID)
  if (!parentPath.length) {
    parentPath.push('me')
  }
  return keywords.some(path => parentPath.includes(path))
}

export interface NavSubItem {
  id: string
  testID: string
  label: string
  link: NavItemLink
}

interface NavItem {
  id: string
  testID: string
  label: string
  shortLabel?: string
  link: NavItemLink
  icon: IconFont
  menu?: NavSubItem[]
  activeKeywords: string[]
}

const generateNavItems = (orgID: string): NavItem[] => {
  const orgPrefix = `/orgs/${orgID}`

  return [
    {
      id: 'data',
      testID: 'nav-item-data',
      icon: IconFont.DisksNav,
      label: 'Data',
      shortLabel: 'Data',
      link: {
        type: 'link',
        location: `${orgPrefix}/data/otcls`,
      },
      activeKeywords: ['data', 'otcls', 'scrapes'],
      menu: [
        {
          id: 'otcls',
          testID: 'nav-subitem-otcls',
          label: 'Otcls',
          link: {
            type: 'link',
            location: `${orgPrefix}/data/otcls`,
          },
        },
        {
          id: 'scrapers',
          testID: 'nav-subitem-scrapers',
          label: 'Scrapers',
          link: {
            type: 'link',
            location: `${orgPrefix}/data/scrapers`,
          },
        },
      ],
    },
    {
      id: 'logs',
      testID: 'nav-item-logs',
      icon: IconFont.Eye,
      label: 'Logs',
      shortLabel: 'Logs',
      link: {
        type: 'link',
        location: `${orgPrefix}/logs`,
      },
      activeKeywords: ['logs'],
    },
    {
      id: 'metrics',
      testID: 'nav-item-metrics',
      icon: IconFont.BarChart,
      label: 'Metrics',
      shortLabel: 'Metrics',
      link: {
        type: 'link',
        location: `${orgPrefix}/metrics`,
      },
      activeKeywords: ['metrics'],
    },
    {
      id: 'traces',
      testID: 'nav-item-traces',
      icon: IconFont.Brush,
      label: 'Traces',
      shortLabel: 'Traces',
      link: {
        type: 'link',
        location: `${orgPrefix}/traces`,
      },
      activeKeywords: ['traces'],
    },
    {
      id: 'profile',
      testID: 'nav-item-profile',
      icon: IconFont.Erlenmeyer,
      label: 'Profile',
      shortLabel: 'Prof',
      link: {
        type: 'link',
        location: `${orgPrefix}/profile`,
      },
      activeKeywords: ['profile'],
    },
    {
      id: 'dashboards',
      testID: 'nav-item-dashboards',
      icon: IconFont.Dashboards,
      label: 'Dashboards',
      shortLabel: 'Boards',
      link: {
        type: 'link',
        location: `${orgPrefix}/dashboards`,
      },
      activeKeywords: ['dashboards'],
    },
    {
      id: 'alerts',
      testID: 'nav-item-alerting',
      icon: IconFont.Bell,
      label: 'Alerts',
      link: {
        type: 'link',
        location: `${orgPrefix}/alerts/checks`,
      },
      activeKeywords: ['alerts'],
      menu: [
        {
          id: 'checks',
          testID: 'nav-subitem-checks',
          label: 'Checks',
          link: {
            type: 'link',
            location: `${orgPrefix}/alerts/checks`,
          },
        },
        {
          id: 'notificationEndpoints',
          testID: 'nav-subitem-checks',
          label: 'notificationEndpoints',
          link: {
            type: 'link',
            location: `${orgPrefix}/alerts/notificationEndpoints`,
          },
        },
      ],
    },
    {
      id: 'settings',
      testID: 'nav-item-settings',
      icon: IconFont.WrenchNav,
      label: 'Settings',
      link: {
        type: 'link',
        location: `${orgPrefix}/settings/variables`,
      },
      activeKeywords: ['settings'],
      menu: [
        {
          id: 'variables',
          testID: 'nav-subitem-variables',
          label: 'Variables',
          link: {
            type: 'link',
            location: `${orgPrefix}/settings/variables`,
          },
        },
        {
          id: 'secrets',
          testID: 'nav-subitem-secrets',
          label: 'Secrets',
          link: {
            type: 'link',
            location: `${orgPrefix}/settings/secrets`,
          },
        },
        {
          id: 'templates',
          testID: 'nav-subitem-templates',
          label: 'Templates',
          link: {
            type: 'link',
            location: `${orgPrefix}/settings/templates`,
          },
        },
        {
          id: 'labels',
          testID: 'nav-subitem-labels',
          label: 'Labels',
          link: {
            type: 'link',
            location: `${orgPrefix}/settings/labels`,
          },
        },
      ],
    },
  ]
}

const Nav: React.FC = () => {
  const orgID = useOrgID()
  const orgPrefix = `/orgs/${orgID}`
  const navItems = generateNavItems(orgID)

  const {inPresentationMode} = usePresentationMode()
  if (inPresentationMode) {
    return null
  }

  return (
    <TreeNav
      expanded={false}
      headerElement={
        <TreeNav.Header
          id="getting-started"
          icon={<Icon glyph={IconFont.CuboNav} />}
          label={<InfluxDBCloudLogo cloud={false} />}
          color={ComponentColor.Secondary}
          linkElement={className => (
            <Link className={className} to={orgPrefix} />
          )}
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
          <TreeNav.Item
            key={item.id}
            id={item.id}
            testID={item.testID}
            icon={<Icon glyph={item.icon} />}
            label={item.label}
            shortLabel={item.shortLabel}
            active={getNavItemActivation(
              item.activeKeywords,
              window.location.pathname
            )}
            linkElement={linkElement}
          >
            {Boolean(item.menu) && (
              <TreeNav.SubMenu>
                {item.menu?.map(menuItem => {
                  const linkElement = (className: string): JSX.Element => {
                    if (menuItem.link.type === 'href') {
                      return (
                        <a
                          href={menuItem.link.location}
                          className={className}
                        />
                      )
                    }

                    return (
                      <Link to={menuItem.link.location} className={className} />
                    )
                  }

                  return (
                    <TreeNav.SubItem
                      key={menuItem.id}
                      id={menuItem.id}
                      testID={menuItem.testID}
                      active={getNavItemActivation(
                        [menuItem.id],
                        window.location.pathname
                      )}
                      label={menuItem.label}
                      linkElement={linkElement}
                    />
                  )
                })}
              </TreeNav.SubMenu>
            )}
          </TreeNav.Item>
        )
      })}
    </TreeNav>
  )
}

export default Nav
