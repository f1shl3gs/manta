import React, {useCallback, useMemo, useState} from 'react'

import {useViewRange} from './useViewRange'
import {useSearch} from './useSearch'

import {useChildrenState} from './useChildrenState'
import {useDetailState} from './useDetailState'
import useHoverIndentGuide from './useHoverIndentGuide'

import {
  ButtonProps,
  Elements,
  Trace,
  TraceKeyValuePair,
  TraceLink,
  TraceSpan,
  UIElementsContext,
  TracePageHeader,
  TraceTimelineViewer,
  transformTraceData,
  TTraceTimeline,
} from './jaeger'
import {Button, Form, Input, Page} from '@influxdata/clockface'

const pageContentsClassName = `alerting-index alerting-index__${'check'}`

const demoData = {
  traceID: '1e7e829e497b015b',
  spans: [
    {
      traceID: '1e7e829e497b015b',
      spanID: '5294bf87375936c1',
      flags: 1,
      operationName: 'bolt.(*KVStore).View',
      references: [
        {
          refType: 'CHILD_OF',
          traceID: '1e7e829e497b015b',
          spanID: '1e7e829e497b015b',
        },
      ],
      startTime: 1618143900000523,
      duration: 6,
      tags: [
        {
          key: 'internal.span.format',
          type: 'string',
          value: 'proto',
        },
      ],
      logs: [
        {
          timestamp: 1618143900000524,
          fields: [
            {
              key: 'filename',
              type: 'string',
              value: '/home/f1shl3gs/Workspaces/goland/manta/bolt/kv.go',
            },
            {
              key: 'line',
              type: 'int64',
              value: 174,
            },
          ],
        },
      ],
      processID: 'p1',
      warnings: null,
    },
    {
      traceID: '1e7e829e497b015b',
      spanID: '18c20f367dc57241',
      flags: 1,
      operationName: 'kv.(*Service).FindTaskByID',
      references: [
        {
          refType: 'CHILD_OF',
          traceID: '1e7e829e497b015b',
          spanID: '1e7e829e497b015b',
        },
      ],
      startTime: 1618143900000521,
      duration: 8,
      tags: [
        {
          key: 'internal.span.format',
          type: 'string',
          value: 'proto',
        },
      ],
      logs: [
        {
          timestamp: 1618143900000522,
          fields: [
            {
              key: 'filename',
              type: 'string',
              value: '/home/f1shl3gs/Workspaces/goland/manta/kv/task.go',
            },
            {
              key: 'line',
              type: 'int64',
              value: 24,
            },
          ],
        },
      ],
      processID: 'p1',
      warnings: null,
    },
    {
      traceID: '1e7e829e497b015b',
      spanID: '2d0ac0850817540f',
      flags: 1,
      operationName: 'bolt.(*KVStore).Update',
      references: [
        {
          refType: 'CHILD_OF',
          traceID: '1e7e829e497b015b',
          spanID: '1e7e829e497b015b',
        },
      ],
      startTime: 1618143900000531,
      duration: 31794,
      tags: [
        {
          key: 'internal.span.format',
          type: 'string',
          value: 'proto',
        },
      ],
      logs: [
        {
          timestamp: 1618143900000532,
          fields: [
            {
              key: 'filename',
              type: 'string',
              value: '/home/f1shl3gs/Workspaces/goland/manta/bolt/kv.go',
            },
            {
              key: 'line',
              type: 'int64',
              value: 197,
            },
          ],
        },
      ],
      processID: 'p1',
      warnings: null,
    },
    {
      traceID: '1e7e829e497b015b',
      spanID: '5d9f1c8d83123e91',
      flags: 1,
      operationName: 'kv.(*Service).AddRunLog',
      references: [
        {
          refType: 'CHILD_OF',
          traceID: '1e7e829e497b015b',
          spanID: '1e7e829e497b015b',
        },
      ],
      startTime: 1618143900032330,
      duration: 36947,
      tags: [
        {
          key: 'internal.span.format',
          type: 'string',
          value: 'proto',
        },
      ],
      logs: [
        {
          timestamp: 1618143900032332,
          fields: [
            {
              key: 'filename',
              type: 'string',
              value:
                '/home/f1shl3gs/Workspaces/goland/manta/kv/task_control.go',
            },
            {
              key: 'line',
              type: 'int64',
              value: 246,
            },
          ],
        },
      ],
      processID: 'p1',
      warnings: null,
    },
    {
      traceID: '1e7e829e497b015b',
      spanID: '33394625a1a9181a',
      flags: 1,
      operationName: 'bolt.(*KVStore).Update',
      references: [
        {
          refType: 'CHILD_OF',
          traceID: '1e7e829e497b015b',
          spanID: '1e7e829e497b015b',
        },
      ],
      startTime: 1618143900069280,
      duration: 30259,
      tags: [
        {
          key: 'internal.span.format',
          type: 'string',
          value: 'proto',
        },
      ],
      logs: [
        {
          timestamp: 1618143900069282,
          fields: [
            {
              key: 'filename',
              type: 'string',
              value: '/home/f1shl3gs/Workspaces/goland/manta/bolt/kv.go',
            },
            {
              key: 'line',
              type: 'int64',
              value: 197,
            },
          ],
        },
      ],
      processID: 'p1',
      warnings: null,
    },
    {
      traceID: '1e7e829e497b015b',
      spanID: '5619693b59390b26',
      flags: 1,
      operationName: 'bolt.(*KVStore).Update',
      references: [
        {
          refType: 'CHILD_OF',
          traceID: '1e7e829e497b015b',
          spanID: '5d9f1c8d83123e91',
        },
      ],
      startTime: 1618143900032333,
      duration: 36944,
      tags: [
        {
          key: 'internal.span.format',
          type: 'string',
          value: 'proto',
        },
      ],
      logs: [
        {
          timestamp: 1618143900032333,
          fields: [
            {
              key: 'filename',
              type: 'string',
              value: '/home/f1shl3gs/Workspaces/goland/manta/bolt/kv.go',
            },
            {
              key: 'line',
              type: 'int64',
              value: 197,
            },
          ],
        },
      ],
      processID: 'p1',
      warnings: null,
    },
    {
      traceID: '1e7e829e497b015b',
      spanID: '1f18132465c81cc9',
      flags: 1,
      operationName: 'kv.(*Service).FindCheckByID',
      references: [
        {
          refType: 'CHILD_OF',
          traceID: '1e7e829e497b015b',
          spanID: '5a3bd43739e9941f',
        },
      ],
      startTime: 1618143900099543,
      duration: 3,
      tags: [
        {
          key: 'internal.span.format',
          type: 'string',
          value: 'proto',
        },
      ],
      logs: [
        {
          timestamp: 1618143900099544,
          fields: [
            {
              key: 'filename',
              type: 'string',
              value: '/home/f1shl3gs/Workspaces/goland/manta/kv/check.go',
            },
            {
              key: 'line',
              type: 'int64',
              value: 21,
            },
          ],
        },
      ],
      processID: 'p1',
      warnings: null,
    },
    {
      traceID: '1e7e829e497b015b',
      spanID: '65953babea29b4ca',
      flags: 1,
      operationName: 'kv.(*Service).AddRunLog',
      references: [
        {
          refType: 'CHILD_OF',
          traceID: '1e7e829e497b015b',
          spanID: '1e7e829e497b015b',
        },
      ],
      startTime: 1618143900099548,
      duration: 30368,
      tags: [
        {
          key: 'internal.span.format',
          type: 'string',
          value: 'proto',
        },
      ],
      logs: [
        {
          timestamp: 1618143900099549,
          fields: [
            {
              key: 'filename',
              type: 'string',
              value:
                '/home/f1shl3gs/Workspaces/goland/manta/kv/task_control.go',
            },
            {
              key: 'line',
              type: 'int64',
              value: 246,
            },
          ],
        },
      ],
      processID: 'p1',
      warnings: null,
    },
    {
      traceID: '1e7e829e497b015b',
      spanID: '63509de5930e61db',
      flags: 1,
      operationName: 'bolt.(*KVStore).Update',
      references: [
        {
          refType: 'CHILD_OF',
          traceID: '1e7e829e497b015b',
          spanID: '65953babea29b4ca',
        },
      ],
      startTime: 1618143900099549,
      duration: 30365,
      tags: [
        {
          key: 'internal.span.format',
          type: 'string',
          value: 'proto',
        },
      ],
      logs: [
        {
          timestamp: 1618143900099550,
          fields: [
            {
              key: 'filename',
              type: 'string',
              value: '/home/f1shl3gs/Workspaces/goland/manta/bolt/kv.go',
            },
            {
              key: 'line',
              type: 'int64',
              value: 197,
            },
          ],
        },
      ],
      processID: 'p1',
      warnings: null,
    },
    {
      traceID: '1e7e829e497b015b',
      spanID: '586defdf7fa0f6c2',
      flags: 1,
      operationName: 'bolt.(*KVStore).Update',
      references: [
        {
          refType: 'CHILD_OF',
          traceID: '1e7e829e497b015b',
          spanID: '1e7e829e497b015b',
        },
      ],
      startTime: 1618143900129920,
      duration: 30277,
      tags: [
        {
          key: 'internal.span.format',
          type: 'string',
          value: 'proto',
        },
      ],
      logs: [
        {
          timestamp: 1618143900129932,
          fields: [
            {
              key: 'filename',
              type: 'string',
              value: '/home/f1shl3gs/Workspaces/goland/manta/bolt/kv.go',
            },
            {
              key: 'line',
              type: 'int64',
              value: 197,
            },
          ],
        },
      ],
      processID: 'p1',
      warnings: null,
    },
    {
      traceID: '1e7e829e497b015b',
      spanID: '5a3bd43739e9941f',
      flags: 1,
      operationName: 'check',
      references: [
        {
          refType: 'CHILD_OF',
          traceID: '1e7e829e497b015b',
          spanID: '1e7e829e497b015b',
        },
      ],
      startTime: 1618143900099542,
      duration: 4,
      tags: [
        {
          key: 'internal.span.format',
          type: 'string',
          value: 'proto',
        },
      ],
      logs: [
        {
          timestamp: 1618143900099542,
          fields: [
            {
              key: 'filename',
              type: 'string',
              value: '/home/f1shl3gs/Workspaces/goland/manta/checks/checks.go',
            },
            {
              key: 'line',
              type: 'int64',
              value: 56,
            },
          ],
        },
      ],
      processID: 'p1',
      warnings: null,
    },
    {
      traceID: '1e7e829e497b015b',
      spanID: '7211897dd481f7a2',
      flags: 1,
      operationName: 'bolt.(*KVStore).View',
      references: [
        {
          refType: 'CHILD_OF',
          traceID: '1e7e829e497b015b',
          spanID: '1f18132465c81cc9',
        },
      ],
      startTime: 1618143900099545,
      duration: 1,
      tags: [
        {
          key: 'internal.span.format',
          type: 'string',
          value: 'proto',
        },
      ],
      logs: [
        {
          timestamp: 1618143900099545,
          fields: [
            {
              key: 'filename',
              type: 'string',
              value: '/home/f1shl3gs/Workspaces/goland/manta/bolt/kv.go',
            },
            {
              key: 'line',
              type: 'int64',
              value: 174,
            },
          ],
        },
      ],
      processID: 'p1',
      warnings: null,
    },
    {
      traceID: '1e7e829e497b015b',
      spanID: '131fa7e09a4c34df',
      flags: 1,
      operationName: 'kv.(*Service).deleteTask',
      references: [
        {
          refType: 'CHILD_OF',
          traceID: '1e7e829e497b015b',
          spanID: '30d892ea519947fb',
        },
      ],
      startTime: 1618143900184468,
      duration: 8,
      tags: [
        {
          key: 'internal.span.format',
          type: 'string',
          value: 'proto',
        },
      ],
      logs: [
        {
          timestamp: 1618143900184470,
          fields: [
            {
              key: 'filename',
              type: 'string',
              value: '/home/f1shl3gs/Workspaces/goland/manta/kv/task.go',
            },
            {
              key: 'line',
              type: 'int64',
              value: 281,
            },
          ],
        },
      ],
      processID: 'p1',
      warnings: null,
    },
    {
      traceID: '1e7e829e497b015b',
      spanID: '38f3925f717b3588',
      flags: 1,
      operationName: 'bolt.(*KVStore).Update',
      references: [
        {
          refType: 'CHILD_OF',
          traceID: '1e7e829e497b015b',
          spanID: '30d892ea519947fb',
        },
      ],
      startTime: 1618143900160207,
      duration: 30360,
      tags: [
        {
          key: 'internal.span.format',
          type: 'string',
          value: 'proto',
        },
      ],
      logs: [
        {
          timestamp: 1618143900160208,
          fields: [
            {
              key: 'filename',
              type: 'string',
              value: '/home/f1shl3gs/Workspaces/goland/manta/bolt/kv.go',
            },
            {
              key: 'line',
              type: 'int64',
              value: 197,
            },
          ],
        },
      ],
      processID: 'p1',
      warnings: null,
    },
    {
      traceID: '1e7e829e497b015b',
      spanID: '1e7e829e497b015b',
      flags: 1,
      operationName: 'Execute',
      references: [],
      startTime: 1618143900000501,
      duration: 220450,
      tags: [
        {
          key: 'sampler.type',
          type: 'string',
          value: 'probabilistic',
        },
        {
          key: 'sampler.param',
          type: 'float64',
          value: 1,
        },
        {
          key: 'internal.span.format',
          type: 'string',
          value: 'proto',
        },
      ],
      logs: [
        {
          timestamp: 1618143900000517,
          fields: [
            {
              key: 'filename',
              type: 'string',
              value:
                '/home/f1shl3gs/Workspaces/goland/manta/task/backend/executor/executor.go',
            },
            {
              key: 'line',
              type: 'int64',
              value: 38,
            },
          ],
        },
        {
          timestamp: 1618143900032328,
          fields: [
            {
              key: 'run_id',
              type: 'string',
              value: '075c38755ea3d000',
            },
            {
              key: 'task_id',
              type: 'string',
              value: '075c3795fd63d000',
            },
          ],
        },
      ],
      processID: 'p1',
      warnings: null,
    },
    {
      traceID: '1e7e829e497b015b',
      spanID: '75175a28a5d6a2b5',
      flags: 1,
      operationName: 'kv.(*Service).FinishRun',
      references: [
        {
          refType: 'CHILD_OF',
          traceID: '1e7e829e497b015b',
          spanID: '1e7e829e497b015b',
        },
      ],
      startTime: 1618143900190572,
      duration: 30378,
      tags: [
        {
          key: 'internal.span.format',
          type: 'string',
          value: 'proto',
        },
      ],
      logs: [
        {
          timestamp: 1618143900190577,
          fields: [
            {
              key: 'filename',
              type: 'string',
              value:
                '/home/f1shl3gs/Workspaces/goland/manta/kv/task_control.go',
            },
            {
              key: 'line',
              type: 'int64',
              value: 128,
            },
          ],
        },
      ],
      processID: 'p1',
      warnings: null,
    },
    {
      traceID: '1e7e829e497b015b',
      spanID: '7bf4bc6b7823900d',
      flags: 1,
      operationName: 'bolt.(*KVStore).Update',
      references: [
        {
          refType: 'CHILD_OF',
          traceID: '1e7e829e497b015b',
          spanID: '75175a28a5d6a2b5',
        },
      ],
      startTime: 1618143900190579,
      duration: 30370,
      tags: [
        {
          key: 'internal.span.format',
          type: 'string',
          value: 'proto',
        },
      ],
      logs: [
        {
          timestamp: 1618143900190579,
          fields: [
            {
              key: 'filename',
              type: 'string',
              value: '/home/f1shl3gs/Workspaces/goland/manta/bolt/kv.go',
            },
            {
              key: 'line',
              type: 'int64',
              value: 197,
            },
          ],
        },
      ],
      processID: 'p1',
      warnings: null,
    },
    {
      traceID: '1e7e829e497b015b',
      spanID: '30d892ea519947fb',
      flags: 1,
      operationName: 'kv.(*Service).UpdateTask',
      references: [
        {
          refType: 'CHILD_OF',
          traceID: '1e7e829e497b015b',
          spanID: '1e7e829e497b015b',
        },
      ],
      startTime: 1618143900160202,
      duration: 30367,
      tags: [
        {
          key: 'internal.span.format',
          type: 'string',
          value: 'proto',
        },
      ],
      logs: [
        {
          timestamp: 1618143900160204,
          fields: [
            {
              key: 'filename',
              type: 'string',
              value: '/home/f1shl3gs/Workspaces/goland/manta/kv/task.go',
            },
            {
              key: 'line',
              type: 'int64',
              value: 223,
            },
          ],
        },
      ],
      processID: 'p1',
      warnings: null,
    },
  ],
  processes: {
    p1: {
      serviceName: 'mantad',
      tags: [
        {
          key: 'client-uuid',
          type: 'string',
          value: 'e25a3596dda1205',
        },
        {
          key: 'hostname',
          type: 'string',
          value: 'localhost.localdomain',
        },
        {
          key: 'ip',
          type: 'string',
          value: '10.32.10.109',
        },
        {
          key: 'jaeger.version',
          type: 'string',
          value: 'Go-2.25.0',
        },
      ],
    },
  },
  warnings: null,
}

function Divider({className}: {className?: string}) {
  return <div />
}

// This needs to be static to prevent remounting on every render.
export const UIElements: Elements = {
  Popover: (() => null as any) as any,
  Tooltip: (() => null as any) as any,
  Icon: (() => null as any) as any,
  Dropdown: (() => null as any) as any,
  Menu: (() => null as any) as any,
  MenuItem: (() => null as any) as any,
  Button({onClick, children, className}: ButtonProps) {
    return (
      <Button
        // variant={'secondary'}

        onClick={() => console.log('oc')}
        className={className}
      >
        {children}
      </Button>
    )
  },
  Divider,
  Input(props) {
    return <Input {...props} />
  },
  InputGroup({children, className, style}) {
    return (
      <span className={className} style={style}>
        {children}
      </span>
    )
  },
}

const TracePage: React.FC = () => {
  const [slim, setSlim] = useState(false)
  const {
    viewRange,
    updateViewRangeTime,
    updateNextViewRangeTime,
  } = useViewRange()

  // @ts-ignore
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
    <Page titleTag={'Trace'}>
      <Page.Header fullWidth={true}>
        <Page.Title title={'Traces'} />
      </Page.Header>
      <Page.Contents
        fullWidth={true}
        scrollable={true}
        className={pageContentsClassName}
      >
        <UIElementsContext.Provider value={UIElements}>
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
        </UIElementsContext.Provider>
      </Page.Contents>
    </Page>
  )
}

export default TracePage
