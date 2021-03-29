// Libraries
import React from 'react'

// Components
import {Controlled as ReactCodeMirror} from 'react-codemirror2'
import {useActiveQuery} from './useQueries'
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

const QueryEditor: React.FC = () => {
  const {activeQuery, onSetText} = useActiveQuery()

  return (
    <div className={'flux-editor'}>
      <div className={'flux-editor--left-panel'}>
        <ReactCodeMirror
          autoFocus={false}
          autoCursor={true}
          value={activeQuery?.text || ''}
          options={options}
          onBeforeChange={(editor, data, value) => onSetText(value)}
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
