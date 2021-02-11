// Libraries
import React, {useCallback, useEffect, useState} from 'react'
import {CachePolicies, useFetch} from 'use-http'

// Components
import ViewSwitcher from 'shared/components/ViewSwitcher'
import EmptyQueryView from 'shared/components/EmptyQueryView'
import ErrorBoundary from 'shared/components/ErrorBoundary'
import {QueriesProvider} from 'components/timeMachine/useQueries'
import {useViewProperties} from 'shared/useViewProperties'
import {transformPromResp} from 'utils/transform'
import {useAutoRefresh} from 'shared/useAutoRefresh'
import {FromFluxResult, fromRows} from '@influxdata/giraffe'

// Hooks
import {useOrgID} from '../../shared/useOrg'

// Utils
import remoteDataState from 'utils/rds'

interface Props {
  cellID?: string
}

const TimeSeries: React.FC<Props> = (props) => {
  const {viewProperties} = useViewProperties()
  const {start, end, step} = useAutoRefresh()
  const {queries} = viewProperties
  const orgID = useOrgID()
  const [errMsg, setErrMsg] = useState<string | undefined>()
  const [result, setResult] = useState<
    Omit<FromFluxResult, 'schema'> | undefined
  >(() => {
    return {
      table: fromRows([]),
      fluxGroupKeyUnion: [],
    }
  })

  const url = `/api/v1/query_range`
  const {get, loading, error} = useFetch(url, {
    cachePolicy: CachePolicies.NO_CACHE,
    onError: ({error}) => {
      setErrMsg(error.message)
    },
  })

  const fetch = useCallback(() => {
    get(
      `?query=${encodeURI(
        queries[0].text
      )}&start=${start}&end=${end}&step=${step}&orgID=${orgID}`
    ).then((resp) => {
      if (resp === undefined) {
        return
      }

      if (resp.status !== 'success') {
        setErrMsg(resp.error)
        return
      }

      setResult(transformPromResp(resp))
    })
  }, [start, end, step, orgID])

  useEffect(() => {
    fetch()
  }, [start, end, step])

  return (
    <ErrorBoundary>
      <QueriesProvider>
        {/*{result ? (
          <ViewSwitcher properties={viewProperties} giraffeResult={result} />
        ) : (
          <EmptyQueryView
            hasResults={result !== undefined}
            loading={remoteDataState(result, error, loading)}
          />
        )}*/}
        <EmptyQueryView
          queries={queries}
          hasResults={result !== undefined}
          loading={remoteDataState(result, error, loading)}
          errorMessage={errMsg}
        >
          <ViewSwitcher properties={viewProperties} giraffeResult={result!} />
        </EmptyQueryView>
      </QueriesProvider>
    </ErrorBoundary>
  )
}

export default TimeSeries
