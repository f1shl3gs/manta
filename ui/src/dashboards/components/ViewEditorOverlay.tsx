// Libraries
import React, {useCallback} from 'react'

// Components
import {Overlay, SpinnerContainer, TechnoSpinner} from '@influxdata/clockface'
import ViewEditorOverlayHeader from './ViewEditorOverlayHeader'
import TimeMachine from 'components/timeMachine/TimeMachine'

// Types
import {Cell, ViewProperties} from 'types/Dashboard'

// Hooks
import {CellProvider, useCell} from './useCell'
import {ViewPropertiesProvider} from 'shared/useViewProperties'

import {ViewOptionProvider} from 'shared/useViewOption'

interface Props {}

const ViewEditorOverlay: React.FC<Props> = (props) => {
  const {cell, remoteDataState} = useCell()

  return (
    <Overlay
      visible={true}
      className={'veo-overlay'}
      onEscape={(visible) => {
        console.log('on escape')
      }}
    >
      <div className={'veo'}>
        <SpinnerContainer
          spinnerComponent={<TechnoSpinner />}
          loading={remoteDataState}
        >
          <ViewPropertiesProvider
            viewProperties={cell?.viewProperties as ViewProperties}
          >
            <ViewEditorOverlayHeader />

            <div className={'veo-contents'}>
              <TimeMachine
                viewProperties={cell?.viewProperties as ViewProperties}
              />
            </div>
          </ViewPropertiesProvider>
        </SpinnerContainer>
      </div>
    </Overlay>
  )
}

const wrapper = () => (
  <CellProvider>
    <ViewOptionProvider>
      <ViewEditorOverlay />
    </ViewOptionProvider>
  </CellProvider>
)

export default wrapper
