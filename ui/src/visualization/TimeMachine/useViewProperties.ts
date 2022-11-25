import {useState} from 'react'
import constate from 'constate'

import {DashboardQuery, ViewProperties} from 'src/types/Dashboard'

interface State {
  viewProperties: ViewProperties
}

const [ViewPropertiesProvider, useViewProperties] = constate(
  (initialState: State) => {
    const [viewProperties, setViewProperties] = useState<ViewProperties>(() => {
      if (initialState.viewProperties === undefined) {
        return {
          type: 'xy',
          xColumn: 'time',
          yColumn: 'value',
          axes: {
            x: {},
            y: {},
          },
          queries: [
            {name: 'query 1', text: '', hidden: false},
          ] as DashboardQuery[],
        } as ViewProperties
      }

      return {
        ...initialState.viewProperties,
      }
    })

    return {
      viewProperties,
      setViewProperties,
    }
  },
  // useViewProperties
  value => value
)

export {ViewPropertiesProvider, useViewProperties}
