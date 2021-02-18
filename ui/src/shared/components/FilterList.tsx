// Libraries
import React, {PropsWithChildren} from 'react'
import {get} from 'lodash'

interface Props<T> {
  list: T[]
  search: string
  searchKeys: string[]
  children: (list: T[]) => any
}

const FilterList: <T>(
  props: Props<T>
) => React.ReactElement<Props<T>> = props => {
  const {list, children} = props

  const filtered = () => {
    return list
  }

  return children(filtered())
}

export default FilterList
