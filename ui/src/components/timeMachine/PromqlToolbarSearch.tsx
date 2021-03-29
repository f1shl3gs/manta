import React, {useEffect, useState} from 'react'
import {IconFont, Input, InputType} from '@influxdata/clockface'
import useDebounce from '../../shared/useDebounce'

interface Props {
  resourceName: string
  onSearch: (s: string) => void
}

const PromqlToolbarSearch: React.FC<Props> = props => {
  const [text, setText] = useState('')
  const {resourceName, onSearch} = props

  // todo: debounce
  const debouncedValue = useDebounce(text, 100)
  useEffect(() => {
    onSearch(debouncedValue)
  }, [debouncedValue, onSearch])

  return (
    <div className={'flux-toolbar--search'}>
      <Input
        type={InputType.Text}
        icon={IconFont.Search}
        placeholder={`Filter ${resourceName}`}
        onChange={ev => {
          setText(ev.target.value)
        }}
        value={text}
      />
    </div>
  )
}

export default PromqlToolbarSearch
