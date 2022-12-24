// Libraries
import React, {FunctionComponent, useEffect} from 'react'
import {connect, ConnectedProps} from 'react-redux'

// Components
import {Overlay} from '@influxdata/clockface'
import CheckOverlayHeader from 'src/checks/components/CheckOverlayHeader'

// Types
import {AppState} from 'src/types/stores'

// Actions
import {resetCheckBuilder, setName} from 'src/checks/actions/builder'
import {createCheckFromBuilder} from 'src/checks/actions/thunks'

// Hooks
import useEscape from 'src/shared/useEscape'

const mstp = (state: AppState) => {
  return {
    name: state.checkBuilder.name,
  }
}

const mdtp = {
  reset: resetCheckBuilder,
  onSetName: setName,
  onSave: createCheckFromBuilder,
}

const connector = connect(mstp, mdtp)
type Props = ConnectedProps<typeof connector>

const NewCheckOverlay: FunctionComponent<Props> = ({
  name,
  reset,
  onSetName,
  onSave,
}) => {
  const onCancel = useEscape()

  useEffect(() => {
    return () => {
      reset()
    }
  }, [reset])

  return (
    <Overlay visible={true} className={'veo-overlay'}>
      <div className={'veo'}>
        <CheckOverlayHeader
          name={name}
          onSetName={onSetName}
          onCancel={onCancel}
          onSave={onSave}
        />
      </div>
    </Overlay>
  )
}

export default connector(NewCheckOverlay)
