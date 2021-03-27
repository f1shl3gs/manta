// Libraries
import React from 'react'

// Components
import {Controlled as ReactCodeMirror} from 'react-codemirror2'
import {useActiveQuery} from './useQueries'
import PromqlEditor from '../promqlEditor/PromqlEditor'
import {Columns, Grid} from '@influxdata/clockface'

// Constants
const options = {
  tabIndex: 1,
  lineNumbers: true,
  autoRefresh: true,
  theme: 'time-machine',
  completeSingle: false,
}

const QueryEditor: React.FC = () => {
  const {activeQuery, onSetText} = useActiveQuery()

  return (
    /*<ReactCodeMirror
      autoScroll={true}
      value={activeQuery?.text || ''}
      options={options}
      onBeforeChange={(editor, data, value) => onSetText(value)}
    />*/

    <PromqlEditor value={activeQuery?.text || ''} onChange={onSetText} />

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
