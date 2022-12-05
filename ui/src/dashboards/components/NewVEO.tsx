// Libraries
import React, {FunctionComponent, useCallback, useState} from 'react'
import {useDispatch} from 'react-redux'
import {useNavigate, useParams} from 'react-router-dom'

// Components
import ViewEditorOverlayHeader from 'src/dashboards/components/ViewEditorOverlayHeader'
import TimeMachine from 'src/timeMachine'
import {createCell} from 'src/cells/actions/thunk'
import {Page} from '@influxdata/clockface'

const NewVEO: FunctionComponent = () => {
  const dispatch = useDispatch()
  const navigate = useNavigate()
  const [name, setName] = useState('')
  const {dashboardID} = useParams()

  const onSubmit = useCallback(() => {
    dispatch(createCell(dashboardID, name))
  }, [dashboardID, name, dispatch])
  const onCancel = () => {
    navigate(-1)
  }

  return (
    <Page>
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
    </Page>
  )
}

export default NewVEO
