// Libraries
import React from 'react'

// Components
import {Controlled as ReactCodeMirror} from 'react-codemirror2'
import {useOtcl} from './useOtcl'

// Constants
const options = {
  tabIndex: 1,
  mode: 'yaml',
  readonly: true,
  lineNumbers: true,
  autoRefresh: true,
  theme: 'material',
  completeSingle: false,
}

const OtclOverlayContent: React.FC = () => {
  const {otcl, onContentChange} = useOtcl()

  return (
    <div className={'veo-contents'}>
      <div className={'time-machine'}>
        <ReactCodeMirror
          autoCursor
          options={options}
          value={otcl.content}
          onBeforeChange={(_, d_, v) => onContentChange(v)}
        />
      </div>
    </div>
  )
}

export default OtclOverlayContent
