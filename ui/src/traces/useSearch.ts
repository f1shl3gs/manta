import {useMemo, useState} from 'react'

import {filterSpans, TraceSpan} from './jaeger'

/**
 * Controls the state of search input that highlights spans if they match the search string.
 * @param spans
 */
export function useSearch(spans?: TraceSpan[]) {
  const [search, setSearch] = useState('')
  const spanFindMatches: Set<string> | undefined | null = useMemo(() => {
    return search && spans ? filterSpans(search, spans) : undefined
  }, [search, spans])

  return {search, setSearch, spanFindMatches}
}
