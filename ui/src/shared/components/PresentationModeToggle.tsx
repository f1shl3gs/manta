// Libraries
import React, {FunctionComponent} from 'react'
import {connect, ConnectedProps} from 'react-redux'

// Components
import {ComponentSize, IconFont, SquareButton} from '@influxdata/clockface'

// Actions
import {enablePresentationMode} from 'src/shared/actions/app'

const mdtp = {
  enablePresentationMode,
}

const connector = connect(null, mdtp)

type Props = ConnectedProps<typeof connector>

const PresentationModeToggle: FunctionComponent<Props> = ({
  enablePresentationMode,
}) => <SquareButton icon={IconFont.ExpandB} onClick={enablePresentationMode} size={ComponentSize.Small} />

export default connector(PresentationModeToggle)
