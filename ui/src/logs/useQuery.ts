import constate from 'constate'
import {useState} from 'react'

const [QueryProvider, useQuery] = constate(
  () => {
    const [query, setQuery] = useState<string>('')
    return {
      query,
      setQuery,
    }
  },
  value => value
)

export {QueryProvider, useQuery}
