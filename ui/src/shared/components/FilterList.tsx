// Utils
import {get} from 'shared/utils/object'

interface Props<T> {
  list: T[]
  search: string
  searchKeys: string[]
  children: (list: T[]) => any
}

function FilterList<T>(props: Props<T>) {
  const {list, search, searchKeys, children} = props

  const filtered = () => {
    return list.filter(item => {
      return (
        undefined !==
        searchKeys.find(sk => {
          const val = get(item, sk)
          if (val === undefined) {
            return false
          }

          return val.toLowerCase().indexOf(search.toLowerCase()) >= 0
        })
      )
    })
  }

  return children(filtered())
}

export default FilterList
