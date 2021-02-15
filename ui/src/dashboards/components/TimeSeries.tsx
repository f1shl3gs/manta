// Libraries
import React, {useState} from 'react'

// Components
import ViewSwitcher from 'shared/components/ViewSwitcher'
import EmptyQueryView from 'shared/components/EmptyQueryView'
import {QueriesProvider} from 'components/timeMachine/useQueries'
import {useViewProperties} from 'shared/useViewProperties'

// Hooks

// Utils
import useQueryResult from './useQueryResult'
import {RemoteDataState} from '@influxdata/clockface'

interface Props {
  cellID?: string
}

const TimeSeries: React.FC<Props> = (props) => {
  const {viewProperties} = useViewProperties()
  const {queries} = viewProperties
  const [errMsg, setErrMsg] = useState<string | undefined>()

  const result = useQueryResult(queries)

  return (
    <QueriesProvider>
      <EmptyQueryView
        queries={queries}
        hasResults={result !== undefined}
        loading={RemoteDataState.Done}
        errorMessage={errMsg}
      >
        <ViewSwitcher properties={viewProperties} giraffeResult={result!} />
      </EmptyQueryView>
    </QueriesProvider>
  )
}

export default TimeSeries
