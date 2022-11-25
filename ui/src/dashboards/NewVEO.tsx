// Libraries
import React, {FunctionComponent} from 'react'

// Components
import {Overlay} from '@influxdata/clockface'
import ViewEditorOverlayHeader from 'src/dashboards/ViewEditorOverlayHeader'
import TimeMachine from 'src/visualization/TimeMachine'

// Hooks
import useEscape from 'src/shared/useEscape'
import {CellProvider, useCell} from 'src/dashboards/useCell'
import {ViewOptionProvider} from 'src/dashboards/components/useViewOption'

const NewVEO: FunctionComponent = () => {
  const {cell} = useCell()

  const onDismiss = useEscape()

  return (
    <Overlay visible={true} className={'veo-overlay'}>
      <div className={'veo'}>
        <ViewEditorOverlayHeader onDismiss={onDismiss} />

        <div className={'veo-contents'}>
          <TimeMachine viewProperties={cell.viewProperties} />
        </div>
      </div>
    </Overlay>
  )
}

export default () => (
  <CellProvider>
    <ViewOptionProvider>
      <NewVEO />
    </ViewOptionProvider>
  </CellProvider>
)