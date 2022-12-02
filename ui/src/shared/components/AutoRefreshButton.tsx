// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {SquareButton, IconFont} from '@influxdata/clockface'

// Hooks
import {useAutoRefresh} from 'src/shared/useAutoRefresh'

const AutoRefreshButton: FunctionComponent = () => {
  const {refresh} = useAutoRefresh()

  return <SquareButton icon={IconFont.Refresh_New} onClick={refresh} />
}

export default AutoRefreshButton
