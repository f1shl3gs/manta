// Libraries
import {useCallback, useState} from 'react'
import constate from 'constate'

// Hooks
import {useViewProperties} from 'shared/useViewProperties'

// Types
import {DashboardQuery, ViewProperties} from 'types/Dashboard'

const [QueriesProvider, useQueries] = constate(
  () => {
    const {viewProperties, setViewProperties} = useViewProperties()
    const {queries = [{text: '', hidden: false}]} = viewProperties
    const [activeIndex, setActiveIndex] = useState(0)

    const setQueries = (queries: DashboardQuery[]) => {
      setViewProperties((prev: ViewProperties) => {
        return {
          ...prev,
          queries,
        }
      })
    }

    return {
      activeIndex,
      setActiveIndex,
      queries,
      setQueries,
    }
  },
  // useQueries
  value => {
    const {queries, setQueries, setActiveIndex, activeIndex} = value
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
      activeIndex,
      queries,
      addQuery,
      removeQuery,
      setActiveIndex,
      onSetText,
    }
  }
)

export {QueriesProvider, useQueries}
