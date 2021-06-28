import React from 'react'
import {IconFont, Input, InputType} from '@influxdata/clockface'

interface Props {
  searchTerm: string
  setSearchTerm: (v: string) => void
}

const PackageSearch: React.FC<Props> = props => {
  const {searchTerm, setSearchTerm} = props

  return (
    <div>
      <Input
        type={InputType.Text}
        icon={IconFont.Search}
        placeholder={'Search packages'}
        onChange={ev => setSearchTerm(ev.target.value)}
        value={searchTerm}
      />
    </div>
  )
}

export default PackageSearch
