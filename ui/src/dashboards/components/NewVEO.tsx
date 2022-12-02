// Libraries
import React, {FunctionComponent, useState} from 'react'

// Components
import {Overlay} from '@influxdata/clockface'
import ViewEditorOverlayHeader from 'src/dashboards/components/ViewEditorOverlayHeader'
import TimeMachine from 'src/timeMachine'
import {useDispatch} from 'react-redux'
import useEscape from 'src/shared/useEscape'
import {useParams} from 'react-router-dom'
import {DEFAULT_VIEWPROPERTIES} from 'src/constants/dashboard'
import {createCell} from 'src/cells/actions/thunk'

const NewVEO: FunctionComponent = () => {
  const [name, setName] = useState('')
  const {dashboardID} = useParams()
  const dispatch = useDispatch()

  const onCancel = useEscape()
  const onSubmit = () => {
    dispatch(createCell(dashboardID, name, DEFAULT_VIEWPROPERTIES))
  }

  return (
    <Overlay visible={true} className={'veo-overlay'}>
      <div className={'veo'}>
        <ViewEditorOverlayHeader
          name={name}
          onSubmit={onSubmit}
          onCancel={onCancel}
          onRename={setName}
        />

        <div className={'veo-contents'}>
          <TimeMachine />
        </div>
      </div>
    </Overlay>
  )
}

export default NewVEO
