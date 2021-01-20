import {useCallback, useState} from 'react'
import constate from 'constate'

import {DashboardQuery, ViewProperties} from 'types/Dashboard'
import {useViewProperties} from 'shared/useViewProperties'

const [QueriesProvider, useQueries, useActiveQuery] = constate(
  () => {
    const {viewProperties, setViewProperties} = useViewProperties()
    const {queries = [{text: '', hidden: false}]} = viewProperties

    const setQueries = (queries: DashboardQuery[]) => {
      // @ts-ignore
      setViewProperties((prev: ViewProperties) => {
        return {
          ...prev,
          queries,
        }
      })
    }

    const [activeIndex, setActiveIndex] = useState(0)

    return {
      activeIndex,
      setActiveIndex,
      queries,
      setQueries,
    }
  },
  // useQueries
  (value) => {
    const {queries, setQueries, setActiveIndex} = value
    const addQuery = useCallback(() => {
      const next = queries.slice().concat({
        hidden: false,
        text: 'new',
        name: `Query ${queries.length + 1}`,
      })

      setQueries(next)
    }, [queries])

    const removeQuery = useCallback(
      (queryIndex: number) => {
        console.log('remove', queryIndex)

        const next = queries.filter((item, index) => index !== queryIndex)
        setActiveIndex(0)
        setQueries(next)
      },
      [queries]
    )

    return {
      queries,
      addQuery,
      removeQuery,
      setActiveIndex,
    }
  },
  // useActiveQuery
  (value) => {
    const {activeIndex, queries, setQueries} = value
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
      [activeIndex, queries]
    )

    return {
      activeIndex,
      onSetText,
      activeQuery: value.queries[activeIndex],
    }
  }
)

export {QueriesProvider, useQueries, useActiveQuery}
