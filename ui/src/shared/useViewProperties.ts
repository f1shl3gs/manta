import {useCallback, useState} from 'react'
import constate from 'constate'

import {DashboardQuery, ViewProperties, XYViewProperties} from 'types/Dashboard'
import {DEFAULT_TIME_FORMAT} from 'constants/timeFormat'

interface State {
  viewProperties: ViewProperties
}

const [ViewPropertiesProvider, useViewProperties, useQueries] = constate(
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
          queries: [] as DashboardQuery[],
        } as ViewProperties
      }

      return {
        ...initialState.viewProperties,
      }
    })

    const sm = useCallback(
      (vp: ViewProperties) => {
        setViewProperties(vp)
      },
      [viewProperties]
    )

    return {
      viewProperties,
      setViewProperties: sm,
    }
  },
  // useViewProperties
  value => value,
  // useQueries
  value => {
    const {
      viewProperties: {queries = []},
    } = value
    // const { queries = [] } = viewProperties;

    return {
      queries,
    }
  }
)

const useLineView = () => {
  const {viewProperties, setViewProperties} = useViewProperties()
  const properties = viewProperties as XYViewProperties

  const onSetXColumn = useCallback(
    (x: string) => {
      setViewProperties({
        ...properties,
        xColumn: x,
      })
    },
    [properties, setViewProperties]
  )

  const onSetYColumn = useCallback(
    (y: string) => {
      setViewProperties({
        ...properties,
        yColumn: y,
      })
    },
    [properties, setViewProperties]
  )

  const onSetTimeFormat = useCallback(
    (timeFormat: string) => {
      setViewProperties({
        ...properties,
        timeFormat,
      })
    },
    [properties, setViewProperties]
  )

  const onSetHoverDimension = useCallback(
    (hoverDimension: 'x' | 'y' | 'xy' | 'auto') => {
      setViewProperties({
        ...properties,
        hoverDimension,
      })
    },
    [properties, setViewProperties]
  )

  const updateYAxis = useCallback(
    (upd: {[key: string]: string}) => {
      // @ts-ignore
      setViewProperties((prev: XYViewProperties) => {
        return {
          ...prev,

          axes: {
            x: prev.axes.x,
            y: {
              ...prev.axes.y,
              ...upd,
            },
          },
        } as XYViewProperties
      })
    },
    [setViewProperties]
  )

  const onSetYAxisLabel = useCallback(
    (label: string) => {
      updateYAxis({label})
    },
    [updateYAxis]
  )

  const onSetYAxisBase = useCallback(
    (base: string) => {
      updateYAxis({base})
    },
    [updateYAxis]
  )

  const onSetYAxisPrefix = useCallback(
    (prefix: string) => {
      updateYAxis({prefix})
    },
    [updateYAxis]
  )

  const onSetYAxisSuffix = useCallback(
    (suffix: string) => {
      updateYAxis({suffix})
    },
    [updateYAxis]
  )

  return {
    onSetXColumn,
    onSetYColumn,
    onSetTimeFormat,
    onSetHoverDimension,
    onSetYAxisLabel,
    onSetYAxisBase,
    onSetYAxisPrefix,
    onSetYAxisSuffix,
    hoverDimension: properties.hoverDimension,
    xColumn: properties.xColumn || 'time',
    yColumn: properties.yColumn || 'value',
    numericColumns: ['time', 'value'] as string[],
    timeFormat: properties.timeFormat || DEFAULT_TIME_FORMAT,
    axes: properties.axes,
  }
}

export {ViewPropertiesProvider, useViewProperties, useLineView}
