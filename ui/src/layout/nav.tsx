// libraries
import React from 'react';
import { Link, useParams } from 'react-router-dom';

// components
import {
  ComponentColor,
  Icon,
  IconFont,
  InfluxDBCloudLogo,
  TreeNav,
} from '@influxdata/clockface';
import { useOrgID } from "../shared/state/organization/organization";

export interface NavItemLink {
  type: 'link' | 'href';
  location: string;
}

const getNavItemActivation = (
  keywords: string[],
  location: string
): boolean => {
  const ignoreOrgAndOrgID = 3;
  const parentPath = location.split('/').slice(ignoreOrgAndOrgID);
  if (!parentPath.length) {
    parentPath.push('me');
  }
  return keywords.some((path) => parentPath.includes(path));
};

export interface NavSubItem {
  id: string;
  testID: string;
  label: string;
  link: NavItemLink;
}

interface NavItem {
  id: string;
  testID: string;
  label: string;
  shortLabel?: string;
  link: NavItemLink;
  icon: IconFont;
  menu?: NavSubItem[];
  activeKeywords: string[];
}

const generateNavItems = (orgID: string): NavItem[] => {
  const orgPrefix = `/orgs/${orgID}`;

  return [
    {
      id: 'otcl',
      testID: 'nav-item-otcls',
      icon: IconFont.Cloud,
      label: 'OTcl',
      shortLabel: 'OTcl',
      link: {
        type: 'link',
        location: `${orgPrefix}/otcls`,
      },
      activeKeywords: ['otcls'],
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
      id: 'load-data',
      testID: 'nav-item-load-data',
      icon: IconFont.DisksNav,
      label: 'Load Data',
      shortLabel: 'Data',
      link: {
        type: 'link',
        location: `${orgPrefix}/load-data/sources`,
      },
      activeKeywords: ['load-data'],
      menu: [
        {
          id: 'sources',
          testID: 'nav-subitem-sources',
          label: 'Sources',
          link: {
            type: 'link',
            location: `${orgPrefix}/load-data/sources`,
          },
        },
        {
          id: 'buckets',
          testID: 'nav-subitem-buckets',
          label: 'Buckets',
          link: {
            type: 'link',
            location: `${orgPrefix}/load-data/buckets`,
          },
        },
        {
          id: 'telegrafs',
          testID: 'nav-subitem-telegrafs',
          label: 'Telegraf',
          link: {
            type: 'link',
            location: `${orgPrefix}/load-data/telegrafs`,
          },
        },
        {
          id: 'scrapers',
          testID: 'nav-subitem-scrapers',
          label: 'Scrapers',
          link: {
            type: 'link',
            location: `${orgPrefix}/load-data/scrapers`,
          },
        },
        {
          id: 'tokens',
          testID: 'nav-subitem-tokens',
          label: 'Tokens',
          link: {
            type: 'link',
            location: `${orgPrefix}/load-data/tokens`,
          },
        },
      ],
    },
    {
      id: 'data-explorer',
      testID: 'nav-item-data-explorer',
      icon: IconFont.GraphLine,
      label: 'Data Explorer',
      shortLabel: 'Explore',
      link: {
        type: 'link',
        location: `${orgPrefix}/data-explorer`,
      },
      activeKeywords: ['data-explorer'],
    },
    {
      id: 'dashboards',
      testID: 'nav-item-dashboards',
      icon: IconFont.Dashboards,
      label: 'Dashboards',
      shortLabel: 'Boards',
      link: {
        type: 'link',
        location: `${orgPrefix}/dashboards-list`,
      },
      activeKeywords: ['dashboards'],
    },
    {
      id: 'tasks',
      testID: 'nav-item-tasks',
      icon: IconFont.Calendar,
      label: 'Tasks',
      link: {
        type: 'link',
        location: `${orgPrefix}/tasks`,
      },
      activeKeywords: ['tasks'],
    },
    {
      id: 'alerting',
      testID: 'nav-item-alerting',
      icon: IconFont.Bell,
      label: 'Alerts',
      link: {
        type: 'link',
        location: `${orgPrefix}/alerting`,
      },
      activeKeywords: ['alerting'],
      menu: [
        {
          id: 'history',
          testID: 'nav-subitem-history',
          label: 'Alert History',
          link: {
            type: 'link',
            location: `${orgPrefix}/alert-history`,
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
  ];
};

// todo: set it manully, test only
const navItems = generateNavItems('06b88c483da3d000');

const Nav: React.FC = () => {
  const orgID = useOrgID();
  const orgPrefix = `/orgs/${orgID}`;

  return (
    <TreeNav
      expanded={false}
      headerElement={
        <TreeNav.Header
          id="getting-started"
          icon={<Icon glyph={IconFont.CuboNav} />}
          label={<InfluxDBCloudLogo cloud={false} />}
          color={ComponentColor.Secondary}
          linkElement={(className) => (
            <Link className={className} to={orgPrefix} />
          )}
        />
      }
    >
      {navItems.map((item) => {
        const linkElement = (classname: string): JSX.Element => {
          if (item.link.type === 'href') {
            return <a href={item.link.location} className={classname} />;
          }

          return <Link to={item.link.location} className={classname} />;
        };

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
                {item.menu?.map((menuItem) => {
                  const linkElement = (className: string): JSX.Element => {
                    if (menuItem.link.type === 'href') {
                      return (
                        <a
                          href={menuItem.link.location}
                          className={className}
                        />
                      );
                    }

                    return (
                      <Link to={menuItem.link.location} className={className} />
                    );
                  };

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
                  );
                })}
              </TreeNav.SubMenu>
            )}
          </TreeNav.Item>
        );
      })}
    </TreeNav>
  );
};

export default Nav;
