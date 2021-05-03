// Libraries
import React, {useCallback} from 'react'

// Components
import {Controlled as ReactCodeMirror} from 'react-codemirror2'
import PromqlToolbar from 'components/timeMachine/PromqlToolbar'

// Constants
const options = {
  tabIndex: 1,
  mode: 'json',
  readonly: true,
  lineNumbers: true,
  autoRefresh: true,
  theme: 'time-machine',
  completeSingle: false,
  scrollbarStyle: 'native',
}

interface Props {
  expr: string
  onExprUpdate: (v: string) => void
  onSubmit: (expr: string) => void
}

const QueryBuilder: React.FC<Props> = props => {
  const {expr, onExprUpdate, onSubmit} = props

  const onKeyPress = useCallback(
    (editor, event: KeyboardEvent) => {
      if (!event.ctrlKey) {
        return
      }

      if (event.key === 'Enter') {
        onSubmit(editor.getValue())
        // onSubmit(exprCopy)
      }
    },
    [onSubmit]
  )

  return (
    <div className={'flux-editor'}>
      <div className={'flux-editor--left-panel'}>
        <ReactCodeMirror
          autoFocus={false}
          autoCursor={true}
          value={expr}
          options={options}
          onBeforeChange={(editor, data, value) => onExprUpdate(value)}
          // @ts-ignore
          onKeyPress={onKeyPress}
        />
      </div>

      <div className={'flux-editor--right-panel'}>
        <PromqlToolbar />
      </div>
    </div>
  )
}

export default QueryBuilder
