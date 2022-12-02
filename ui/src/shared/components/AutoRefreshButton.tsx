// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {SquareButton, IconFont} from '@influxdata/clockface'

const AutoRefreshButton: FunctionComponent = () => {
  return <SquareButton icon={IconFont.Refresh_New} onClick={refresh} />
}

export default AutoRefreshButton
