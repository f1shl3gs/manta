import React, {FunctionComponent} from 'react'

interface Props {
  message: string
  testID?: string
}

const EmptyGraphMessage: FunctionComponent<Props> = ({message, testID}) => (
  <div className={'cell--view-empty'} data-testid={testID}>
    <h4>{message}</h4>
  </div>
)

export default EmptyGraphMessage
