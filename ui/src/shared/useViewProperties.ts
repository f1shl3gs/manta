import {useCallback, useState} from 'react'
import constate from 'constate'

import {ViewProperties, XYViewProperties} from 'types/Dashboard'
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
          queries: [
            {
              text: '',
              hidden: false,
            },
          ],
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
  (value) => value,
  // useQueries
  (value) => {
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
    [properties]
  )

  const onSetYColumn = useCallback(
    (y: string) => {
      setViewProperties({
        ...properties,
        yColumn: y,
      })
    },
    [properties]
  )

  const onSetTimeFormat = useCallback(
    (timeFormat: string) => {
      setViewProperties({
        ...properties,
        timeFormat,
      })
    },
    [properties]
  )

  const onSetHoverDimension = useCallback(
    (hoverDimension: 'x' | 'y' | 'xy' | 'auto') => {
      console.log('on hover set', properties)

      setViewProperties({
        ...properties,
        hoverDimension,
      })
    },
    [properties]
  )

  const updateYAxis = useCallback((upd: {[key: string]: string}) => {
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
  }, [])

  const onSetYAxisLabel = useCallback(
    (label: string) => {
      /*setViewProperties({
      ...properties,
      axes: {
        x: properties.axes.x,
        y: {
          ...properties.axes.y,
          label
        }
      }
    });*/
      updateYAxis({label})
    },
    [properties]
  )

  const onSetYAxisBase = useCallback(
    (base: string) => {
      console.log(properties)
      updateYAxis({base})
      /*setViewProperties({
      ...properties,
      axes: {
        x: properties.axes.x,
        y: {
          ...properties.axes.y,
          base
        }
      }
    });*/
    },
    [properties]
  )

  const onSetYAxisPrefix = useCallback(
    (prefix: string) => {
      updateYAxis({prefix})
      /*setViewProperties({
      ...properties,
      axes: {
        x: properties.axes.x,
        y: {
          ...properties.axes.y,
          prefix
        }
      }
    });*/
    },
    [properties]
  )

  const onSetYAxisSuffix = useCallback(
    (suffix: string) => {
      updateYAxis({suffix})
      /*setViewProperties({
      ...properties,
      axes: {
        x: properties.axes.x,
        y: {
          ...properties.axes.y,
          suffix
        }
      }
    });*/
    },
    [properties]
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
