import {useCallback, useState} from 'react'
import constate from 'constate'

import {DashboardQuery, ViewProperties} from 'src/types/Dashboard'

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
            {name: 'query 1', text: '', hidden: false},
          ] as DashboardQuery[],
        } as ViewProperties
      }

      return {
        ...initialState.viewProperties,
      }
    })

    const setQueries = (queries: DashboardQuery[]) => {
      setViewProperties((prev: ViewProperties) => {
        return {
          ...prev,
          queries,
        }
      })
    }

    return {
      viewProperties,
      setViewProperties,
      setQueries,
    }
  },
  // useViewProperties
  value => value,
  // useQueries
  value => {
    const setQueries = value.setQueries
    const {queries} = value.viewProperties
    const [activeIndex, setActiveIndex] = useState(0)

    const addQuery = useCallback(() => {
      const next = queries.slice().concat({
        hidden: false,
        text: 'new',
        name: `Query ${queries.length + 1}`,
      })

      setQueries(next)
    }, [queries, setQueries])

    const removeQuery = useCallback(
      (queryIndex: number) => {
        if (queries.length === 1) {
          // nothing to delete anymore
          return
        }

        console.log('remove', queryIndex)

        const next = queries.filter((item, index) => index !== queryIndex)
        console.log('set index', 0)
        setActiveIndex(0)
        console.log('set index done')

        console.log('set queries', next)
        setQueries(next)
        console.log('set queries done')
      },
      [queries, setActiveIndex, setQueries]
    )

    const onSetText = useCallback(
      (text: string) => {
        const next = queries.map((item, queryIndex) => {
          if (queryIndex !== activeIndex) {
            return item
          }

          return {
            ...item,
            text,
          }
        })

        setQueries(next)
      },
      [activeIndex, queries, setQueries]
    )

    return {
      addQuery,
      removeQuery,
      onSetText,
      queries,
      activeQuery: queries[activeIndex],
      activeIndex,
      setActiveIndex,
      setQueries: value.setQueries,
    }
  }
)

export {ViewPropertiesProvider, useViewProperties, useQueries}
