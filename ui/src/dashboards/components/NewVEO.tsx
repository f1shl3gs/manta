// Libraries
import React, {FunctionComponent, useCallback, useEffect, useState} from 'react'
import {useDispatch} from 'react-redux'
import {useNavigate, useParams} from 'react-router-dom'

// Components
import ViewEditorOverlayHeader from 'src/dashboards/components/ViewEditorOverlayHeader'
import TimeMachine from 'src/timeMachine/components/TimeMachine'
import {Page} from '@influxdata/clockface'

// Actions
import {createCell} from 'src/cells/actions/thunk'
import {resetTimeMachine} from 'src/timeMachine/actions'

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

  useEffect(() => {
    return () => {
      dispatch(resetTimeMachine())
    }
  }, [dispatch])

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
