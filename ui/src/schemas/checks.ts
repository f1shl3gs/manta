// Libraries
import {schema} from 'normalizr'

// Types
import {Check, CheckBase, Conditions} from 'src/types/checks'
import {ResourceType} from 'src/types/resources'
import {RemoteDataState} from '@influxdata/clockface'

export const checkSchema = new schema.Entity(
  ResourceType.Checks,
  {},
  {
    processStrategy: (base: CheckBase) => {
      const conditions: Conditions = {}
      base.conditions.forEach(c => {
        conditions[c.status] = c
      })

      const check: Check = {
        ...base,
        conditions,
        activeStatus: base.status,
        status: RemoteDataState.Done,
      }

      return check
    },
  }
)

export const arrayOfChecks = [checkSchema]
