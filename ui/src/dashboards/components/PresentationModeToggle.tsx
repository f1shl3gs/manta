// Libraries
import React, {FunctionComponent} from 'react'
import {connect, ConnectedProps} from 'react-redux'

// Components
import {IconFont, SquareButton} from '@influxdata/clockface'

// Actions
import {enablePresentationMode} from 'src/shared/actions/app'

type Props = ConnectedProps<typeof connector>

const PresentationModeToggle: FunctionComponent<Props> = ({enablePresentationMode}) => (
  <SquareButton icon={IconFont.ExpandB} onClick={enablePresentationMode} />
  )

const mdtp = {
  enablePresentationMode: enablePresentationMode
}

const connector = connect(null, mdtp)

export default connector(PresentationModeToggle)
