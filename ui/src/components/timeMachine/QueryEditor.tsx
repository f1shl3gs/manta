// Libraries
import React from 'react'

// Components
import {Controlled as ReactCodeMirror} from 'react-codemirror2'
import PromqlToolbar from './PromqlToolbar'

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
  query: string
  onChange: (v: string) => void
}

// todo: maybe rename it to PromqlEditor

const QueryEditor: React.FC<Props> = ({query, onChange}) => {
  return (
    <div className={'flux-editor'}>
      <div className={'flux-editor--left-panel'}>
        <ReactCodeMirror
          autoFocus={false}
          autoCursor={true}
          value={query}
          options={options}
          onBeforeChange={(editor, data, value) => onChange(value)}
        />
      </div>

      <div className={'flux-editor--right-panel'}>
        <PromqlToolbar />
      </div>
    </div>

    /*<PromqlEditor value={activeQuery?.text || ''} onChange={onSetText} />*/

    /*<Grid style={{height: '100%'}}>
      <Grid.Row>
        <Grid.Column widthSM={Columns.Six}>
          <PromqlEditor
            value={activeQuery?.text || ''}
            onChange={v => console.log('onchange', v)}
          />
        </Grid.Column>

        <Grid.Column widthSM={Columns.Six}>
          <ReactCodeMirror
            autoScroll={true}
            value={activeQuery?.text || ''}
            options={options}
            onBeforeChange={(editor, data, value) => onSetText(value)}
          />
        </Grid.Column>
      </Grid.Row>
    </Grid>*/
  )
}

export default QueryEditor
