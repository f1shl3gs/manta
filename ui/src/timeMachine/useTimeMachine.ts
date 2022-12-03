import constate from 'constate'
import {useCallback, useMemo, useState} from 'react'
import {DEFAULT_VIEWPROPERTIES} from 'src/constants/dashboard'
import {ViewProperties} from 'src/types/dashboard'

interface State {
  viewProperties?: ViewProperties
}

const [TimeMachineProvider, useTimeMachine, useViewingVisOptions, useQueries] =
  constate(
    (state: State) => {
      const [viewProperties, setViewProperties] = useState<ViewProperties>(
        state.viewProperties ?? DEFAULT_VIEWPROPERTIES
      )
      const [viewingVisOptions, setViewingVisOptions] = useState(false)
      const queries = viewProperties.queries
      const [activeIndex, setActiveIndex] = useState(0)

      const addQuery = useCallback(() => {
        const next = queries.slice().concat({
          hidden: false,
          text: '',
          name: `Query ${queries.length + 1}`,
        })

        setViewProperties(prev => ({
          ...prev,
          queries: next,
        }))
      }, [queries, setViewProperties])

      const removeQuery = useCallback(
        (queryIndex: number) => {
          if (queries.length === 1) {
            // nothing to delete anymore
            return
          }

          const next = queries.filter((item, index) => index !== queryIndex)
          setActiveIndex(0)

          setViewProperties(prev => ({
            ...prev,
            queries: next,
          }))
        },
        [queries, setActiveIndex, setViewProperties]
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

          setViewProperties(prev => ({
            ...prev,
            queries: next,
          }))
        },
        [activeIndex, queries, setViewProperties]
      )

      return {
        viewProperties,
        setViewProperties,
        viewingVisOptions,
        setViewingVisOptions,
        queries,
        activeIndex,
        setActiveIndex,
        addQuery,
        onSetText,
        removeQuery,
        activeQuery: queries[activeIndex],
      }
    },

    // useTimeMachine
    value =>
      useMemo(() => {
        return {
          ...value,
        }
      }, [value]),

    // useViewingVisOptions
    value =>
      useMemo(
        () => ({
          viewingVisOptions: value.viewingVisOptions,
          setViewingVisOptions: value.setViewingVisOptions,
        }),
        [value.viewingVisOptions, value.setViewingVisOptions]
      ),
    // useQueries
    value =>
      useMemo(
        () => ({
          queries: value.queries,
          activeIndex: value.activeIndex,
          setActiveIndex: value.setActiveIndex,
          addQuery: value.addQuery,
          onSetText: value.onSetText,
          removeQuery: value.removeQuery,
          activeQuery: value.queries[value.activeIndex],
        }),
        [
          value.queries,
          value.activeIndex,
          value.addQuery,
          value.onSetText,
          value.removeQuery,
          value.setActiveIndex,
        ]
      )
  )

export {TimeMachineProvider, useTimeMachine, useViewingVisOptions, useQueries}
