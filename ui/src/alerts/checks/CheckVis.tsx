// Libraries
import React, {useMemo} from 'react'

// Components
import {Plot} from '@influxdata/giraffe'
import XYPlot from '../../components/timeMachine/XYPlot'

// Types
import {DashboardQuery, XYViewProperties} from '../../types/Dashboard'
import useQueryResult from '../../dashboards/components/useQueryResult'
import {ViewPropertiesProvider} from '../../shared/useViewProperties'

interface Props {
  query: string
}

const viewProperties: XYViewProperties = {
  type: 'xy',
  geom: 'line',
  position: 'overlaid',
  queries: [],
  axes: {
    x: {},
    y: {},
  },
}

const CheckVis: React.FC<Props> = ({query}) => {
  const queries = useMemo<DashboardQuery[]>(() => {
    return [
      {
        text: query,
        hidden: false,
      },
    ]
  }, [query])
  const {table, fluxGroupKeyUnion} = useQueryResult(queries)

  return (
    <div className={'time-machine--view'}>
      <ViewPropertiesProvider viewProperties={viewProperties}>
        <XYPlot table={table} groupKeyUnion={fluxGroupKeyUnion}>
          {config => <Plot config={config} />}
        </XYPlot>
      </ViewPropertiesProvider>
    </div>
  )
}

export default CheckVis
