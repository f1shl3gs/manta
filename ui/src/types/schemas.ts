import {Dashboard} from 'src/types/dashboards'
import {Organization} from 'src/types/organization'
import {Cell} from 'src/types/cells'
import {Configuration} from 'src/types/configuration'
import {Scrape} from 'src/types/scrape'

// DashboardEntities defines the result of normalizr's normalization of the
// "dashboards" resource
export interface DashboardEntities {
  dashboards: {
    [uuid: string]: Dashboard
  }
  cells: {
    [uuid: string]: Cell
  }
}

// OrgEntities defines the result of normalizr's normalization
// of the "organizations" resource
export interface OrgEntities {
  orgs: {
    [uuid: string]: Organization
  }
}

// ConfigurationEntities defines the result of normalizr's normalization
// of the `configurations` resource
export interface ConfigurationEntities {
  configurations: {
    [uuid: string]: Configuration
  }
}

export interface ScrapeEntities {
  scrapes: {
    [uuid: string]: Scrape
  }
}
