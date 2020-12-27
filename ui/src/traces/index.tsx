import React, {useCallback, useMemo, useState} from 'react'

import {useViewRange} from './useViewRange'
import {useSearch} from './useSearch'

import {useChildrenState} from './useChildrenState'
import {useDetailState} from './useDetailState'
import useHoverIndentGuide from './useHoverIndentGuide'

import {
  Trace,
  TraceKeyValuePair,
  TraceLink,
  TraceSpan,
} from '../containers/jaeger'
import {
  TracePageHeader,
  TraceTimelineViewer,
  transformTraceData,
  TTraceTimeline,
} from '../containers/jaeger'
import {Page} from '@influxdata/clockface'

const pageContentsClassName = `alerting-index alerting-index__${'check'}`

const demoData = {
  traceID: '66073a21f6927234',
  spans: [
    {
      traceID: '66073a21f6927234',
      spanID: '66073a21f6927234',
      flags: 1,
      operationName: '/api/services/{service}/operations',
      references: [],
      startTime: 1606838245318873,
      duration: 120,
      tags: [
        {key: 'sampler.type', type: 'string', value: 'const'},
        {
          key: 'sampler.param',
          type: 'bool',
          value: true,
        },
        {key: 'span.kind', type: 'string', value: 'server'},
        {
          key: 'http.method',
          type: 'string',
          value: 'GET',
        },
        {
          key: 'http.url',
          type: 'string',
          value: '/api/services/jaeger-query/operations',
        },
        {
          key: 'component',
          type: 'string',
          value: 'net/http',
        },
        {key: 'http.status_code', type: 'int64', value: 200},
        {
          key: 'internal.span.format',
          type: 'string',
          value: 'proto',
        },
      ],
      logs: [],
      processID: 'p1',
      warnings: null,
    },
  ],
  processes: {
    p1: {
      serviceName: 'jaeger-query',
      tags: [
        {key: 'client-uuid', type: 'string', value: '7ad528e472c1870e'},
        {
          key: 'hostname',
          type: 'string',
          value: 'd312f61ac3fc',
        },
        {key: 'ip', type: 'string', value: '10.0.2.100'},
        {
          key: 'jaeger.version',
          type: 'string',
          value: 'Go-2.23.1',
        },
      ],
    },
  },
  warnings: null,
}

const TracePage: React.FC = () => {
  const [slim, setSlim] = useState(false)
  const {
    viewRange,
    updateViewRangeTime,
    updateNextViewRangeTime,
  } = useViewRange()

  const traceProp = transformTraceData(demoData) as Trace
  const {search, setSearch, spanFindMatches} = useSearch(traceProp?.spans)

  const {
    expandOne,
    collapseOne,
    childrenToggle,
    collapseAll,
    childrenHiddenIDs,
    expandAll,
  } = useChildrenState()
  const {
    removeHoverIndentGuideId,
    addHoverIndentGuideId,
    hoverIndentGuideIds,
  } = useHoverIndentGuide()
  const {
    detailStates,
    toggleDetail,
    detailLogItemToggle,
    detailLogsToggle,
    detailProcessToggle,
    detailReferencesToggle,
    detailTagsToggle,
    detailWarningsToggle,
    detailStackTracesToggle,
  } = useDetailState()

  const [spanNameColumnWidth, setSpanNameColumnWidth] = useState(0.25)

  const traceTimeline: TTraceTimeline = useMemo(
    () => ({
      childrenHiddenIDs,
      detailStates,
      hoverIndentGuideIds,
      shouldScrollToFirstUiFindMatch: false,
      spanNameColumnWidth,
      traceID: traceProp?.traceID,
    }),
    [
      childrenHiddenIDs,
      detailStates,
      hoverIndentGuideIds,
      spanNameColumnWidth,
      traceProp?.traceID,
    ]
  )

  // const createSpanLink = useMemo(() => createSpanLinkFactory(splitOpenFn), [splitOpenFn]);
  const createSpanLink = (span: TraceSpan) => {
    /*
    { href: string; onClick?: (e: React.MouseEvent) => void; content: React.ReactNode }
    * */
    return {
      href: `/${span.traceID}`,
      // onClick: (ev: React.MouseEvent) => console.log('span link onclick', ev),
      content: <div>span</div>,
    }
  }

  return (
    <>
      <Page titleTag={'Trace'}>
        <Page.Header fullWidth={true}>
          <Page.Title title={'Traces'}></Page.Title>
        </Page.Header>
        <Page.Contents
          fullWidth={true}
          scrollable={false}
          className={pageContentsClassName}
        >
          <TracePageHeader
            canCollapse={false}
            clearSearch={useCallback(() => {}, [])}
            focusUiFindMatches={useCallback(() => {}, [])}
            hideMap={false}
            hideSummary={false}
            nextResult={useCallback(() => {}, [])}
            onSlimViewClicked={useCallback(() => setSlim(!slim), [])}
            onTraceGraphViewClicked={useCallback(() => {}, [])}
            prevResult={useCallback(() => {}, [])}
            resultCount={0}
            slimView={slim}
            textFilter={null}
            trace={traceProp}
            traceGraphView={false}
            updateNextViewRangeTime={updateNextViewRangeTime}
            updateViewRangeTime={updateViewRangeTime}
            viewRange={viewRange}
            searchValue={search}
            onSearchValueChange={setSearch}
            hideSearchButtons={true}
          />

          <TraceTimelineViewer
            registerAccessors={useCallback(() => {}, [])}
            scrollToFirstVisibleSpan={useCallback(() => {}, [])}
            findMatchesIDs={spanFindMatches}
            trace={traceProp}
            traceTimeline={traceTimeline}
            updateNextViewRangeTime={updateNextViewRangeTime}
            updateViewRangeTime={updateViewRangeTime}
            viewRange={viewRange}
            focusSpan={useCallback(() => {}, [])}
            createLinkToExternalSpan={useCallback(() => '', [])}
            setSpanNameColumnWidth={setSpanNameColumnWidth}
            collapseAll={collapseAll}
            collapseOne={collapseOne}
            expandAll={expandAll}
            expandOne={expandOne}
            childrenToggle={childrenToggle}
            clearShouldScrollToFirstUiFindMatch={useCallback(() => {}, [])}
            detailLogItemToggle={detailLogItemToggle}
            detailLogsToggle={detailLogsToggle}
            detailWarningsToggle={detailWarningsToggle}
            detailStackTracesToggle={detailStackTracesToggle}
            detailReferencesToggle={detailReferencesToggle}
            detailProcessToggle={detailProcessToggle}
            detailTagsToggle={detailTagsToggle}
            detailToggle={toggleDetail}
            setTrace={useCallback(
              (
                trace: Trace | null | undefined,
                uiFind: string | null | undefined
              ) => {},
              []
            )}
            addHoverIndentGuideId={addHoverIndentGuideId}
            removeHoverIndentGuideId={removeHoverIndentGuideId}
            linksGetter={useCallback(
              (
                span: TraceSpan,
                items: TraceKeyValuePair[],
                itemIndex: number
              ) => [] as TraceLink[],
              []
            )}
            uiFind={search}
            // @ts-ignore
            createSpanLink={createSpanLink}
          />
        </Page.Contents>
      </Page>
    </>
  )
}

export default TracePage
