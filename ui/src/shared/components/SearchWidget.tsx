// Libraries
import React from 'react'

// Components
import {IconFont, Input} from '@influxdata/clockface'

interface Props {
  search: string
  placeholder: string
  onSearch: (k: string) => void
}

const SearchWidget: React.FC<Props> = props => {
  const {search, placeholder, onSearch} = props

  return (
    <Input
      icon={IconFont.Search_New}
      placeholder={placeholder}
      value={search}
      onChange={ev => onSearch(ev.target.value)}
      // @ts-ignore
      onBlur={ev => onSearch(ev.target.value)}
      className={'search-widget-input'}
    />
  )
}

export default SearchWidget
