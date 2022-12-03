// Libraries
import React, {FunctionComponent, useCallback, useState} from 'react'

// Components
import ViewEditorOverlayHeader from 'src/dashboards/ViewEditorOverlayHeader'
import TimeMachine from 'src/timeMachine'

// Hooks
import {CellProvider, useCell} from 'src/cells/useCell'
import {TimeMachineProvider} from 'src/timeMachine/useTimeMachine'

const NewVEO: FunctionComponent = () => {
  const {
    cell: {viewProperties},
    createCell,
  } = useCell()
  const [name, setName] = useState('')

  const handleSubmit = useCallback(() => {
    createCell({
      name,
      viewProperties,
      w: 4,
      h: 4,
    })
  }, [createCell, name, viewProperties])

  return (
    <div className={'veo'}>
        <TimeMachineProvider viewProperties={viewProperties}>
          <ViewEditorOverlayHeader
            name={name}
            onRename={setName}
            onSubmit={handleSubmit}
          />

          <div className={'veo-contents'}>
            <TimeMachine />
          </div>
        </TimeMachineProvider>
    </div>
  )
}

export default () => (
  <CellProvider>
    <NewVEO />
  </CellProvider>
)
