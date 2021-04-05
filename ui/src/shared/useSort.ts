import constate from 'constate'
import {useState} from 'react'
import {Sort} from '@influxdata/clockface'

import {SortTypes} from '../types/sort'

const [SortProvider, useSort] = constate(() => {
  const [sortOption, setSortOption] = useState({
    key: 'update',
    type: SortTypes.Date,
    direction: Sort.Descending,
  })

  return {
    ...sortOption,
    setSortOption,
  }
})

export {SortProvider, useSort}
