// Libraries
import React from 'react'

// Components
import {Overlay, SpinnerContainer, TechnoSpinner} from '@influxdata/clockface'
import ViewEditorOverlayHeader from './ViewEditorOverlayHeader'
import TimeMachine from 'components/timeMachine/TimeMachine'

// Types
import {ViewProperties} from 'types/Dashboard'

// Hooks
import {CellProvider, useCell} from './useCell'
import {ViewOptionProvider} from 'dashboards/components/useViewOption'
import {ViewPropertiesProvider} from 'shared/useViewProperties'

interface Props {}

const ViewEditorOverlay: React.FC<Props> = () => {
  const {cell, loading} = useCell()

  return (
    <Overlay
      visible={true}
      className={'veo-overlay'}
      onEscape={visible => {
        console.log('on escape', visible)
      }}
    >
      <div className={'veo'}>
        <SpinnerContainer
          spinnerComponent={<TechnoSpinner />}
          loading={loading}
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
