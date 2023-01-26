import {Dashboard} from 'src/types/dashboards'
import {Organization} from 'src/types/organization'
import {Cell} from 'src/types/cells'
import {Config} from 'src/types/config'
import {Scrape} from 'src/types/scrape'
import {Check} from 'src/types/checks'
import {Secret} from 'src/types/secrets'
import {NotificationEndpoint} from 'src/types/notificationEndpoints'

// ConfigEntities defines the result of normalizr's normalization
// of the `configs` resource
export interface ConfigEntities {
  configs: {
    [uuid: string]: Config
  }
}

export interface CheckEntities {
  checks: {
    [uuid: string]: Check
  }
}

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

export interface NotificationEndpointEntities {
  notificationEndpoints: {
    [uuid: string]: NotificationEndpoint
  }
}

// OrgEntities defines the result of normalizr's normalization
// of the "organizations" resource
export interface OrgEntities {
  orgs: {
    [uuid: string]: Organization
  }
}

export interface SecretEntities {
  secrets: {
    [key: string]: Secret
  }
}

export interface ScrapeEntities {
  scrapes: {
    [uuid: string]: Scrape
  }
}
