import React from 'react'
import PromqlEditor from '../../../components/promqlEditor/PromqlEditor'
import {useCheck} from '../useCheck'
import {Controlled as ReactCodeMirror} from 'react-codemirror2'
import PromqlToolbar from '../../../components/timeMachine/PromqlToolbar'

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

const QueryBuilder: React.FC = () => {
  const {expr, onExprUpdate} = useCheck()

  console.log(expr)

  return (
    <div className={'flux-editor'}>
      <div className={'flux-editor--left-panel'}>
        <ReactCodeMirror
          autoFocus={false}
          autoCursor={true}
          value={expr}
          options={options}
          onBeforeChange={(editor, data, value) => onExprUpdate(value)}
        />
      </div>

      <div className={'flux-editor--right-panel'}>
        <PromqlToolbar />
      </div>
    </div>
  )
}

export default QueryBuilder
