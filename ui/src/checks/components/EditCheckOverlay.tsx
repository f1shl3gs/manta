// Libraries
import React, {FunctionComponent, useEffect} from 'react'
import {connect, ConnectedProps} from 'react-redux'
import {useParams} from 'react-router-dom'

// Components
import {Overlay, SpinnerContainer, TechnoSpinner} from '@influxdata/clockface'
import CheckOverlayHeader from 'src/checks/components/CheckOverlayHeader'
import CheckBuilder from 'src/checks/components/checkBuilder/CheckBuilder'

// Types
import {AppState} from 'src/types/stores'

// Hooks
import useEscape from 'src/shared/useEscape'

// Actions
import {setName, resetCheckBuilder} from 'src/checks/actions/builder'
import {
  getCheckForBuilder,
  createCheckFromBuilder,
} from 'src/checks/actions/thunks'

const mstp = (state: AppState) => {
  const {name, status} = state.checkBuilder

  return {
    name,
    loading: status,
  }
}

const mdtp = {
  getCheckForBuilder,
  resetCheckBuilder,
  onSetName: setName,
  onSave: createCheckFromBuilder,
}

const connector = connect(mstp, mdtp)
type Props = ConnectedProps<typeof connector>

const EditCheckOverlay: FunctionComponent<Props> = ({
  name,
  loading,
  onSetName,
  onSave,
  getCheckForBuilder,
  resetCheckBuilder,
}) => {
  const {id} = useParams<{id: string}>()
  const onCancel = useEscape()

  useEffect(() => {
    getCheckForBuilder(id)

    return () => {
      resetCheckBuilder()
    }
  }, [id, getCheckForBuilder, resetCheckBuilder])

  return (
    <Overlay visible={true}>
      <div className={'veo'}>
        <SpinnerContainer
          loading={loading}
          spinnerComponent={<TechnoSpinner />}
        >
          <CheckOverlayHeader
            name={name}
            onSetName={onSetName}
            onCancel={onCancel}
            onSave={onSave}
          />

          <div className={'veo-contents'}>
            <CheckBuilder />
          </div>
        </SpinnerContainer>
      </div>
    </Overlay>
  )
}

export default connector(EditCheckOverlay)
