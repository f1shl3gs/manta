import React, {FunctionComponent} from 'react'
import {DapperScrollbars, Icon, IconFont} from '@influxdata/clockface'

interface Props {
  message: string
  testID?: string
}

const EmptyGraphError: FunctionComponent<Props> = ({message, testID}) => {
  return (
    <div className={'cell--view-empty'} data-testid={testID}>
      <div className={'empty-graph-error'} data-testid="empty-graph-error">
        <DapperScrollbars
          className={'empty-graph-error--scroll'}
          autoHide={true}
        >
          <pre>
            <Icon glyph={IconFont.AlertTriangle} />
            <code className={'cell--error-message'}>{message}</code>
          </pre>
        </DapperScrollbars>
      </div>
    </div>
  )
}

export default EmptyGraphError
