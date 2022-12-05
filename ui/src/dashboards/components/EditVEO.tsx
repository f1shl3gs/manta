// Libraries
import React, {FunctionComponent, useCallback, useState} from 'react'

// Components
import ViewEditorOverlayHeader from 'src/dashboards/components/ViewEditorOverlayHeader'
import TimeMachine from 'src/timeMachine/components/TimeMachine'
import {Page} from '@influxdata/clockface'

// Hooks
import {useDispatch, useSelector} from 'react-redux'
import {useNavigate, useParams} from 'react-router-dom'

// Types
import {AppState} from 'src/types/stores'
import {Cell} from 'src/types/cells'
import {ResourceType} from 'src/types/resources'

// Actions
import {updateCell} from 'src/cells/actions/thunk'
import GetResource from 'src/resources/components/GetResource'

// Selectors
import {getByID} from 'src/resources/selectors'

const EditVEO: FunctionComponent = () => {
  const dispatch = useDispatch()
  const {dashboardID, cellID} = useParams()
  const navigate = useNavigate()
  const cell = useSelector((state: AppState) => {
    return getByID<Cell>(state, ResourceType.Cells, cellID)
  })
  const [name, setName] = useState(cell.name)
  const viewProperties = useSelector((state: AppState) => {
    return state.timeMachine.viewProperties
  })

  const onCancel = () => {
    navigate(-1)
  }
  const onSubmit = useCallback(() => {
    dispatch(
      updateCell(dashboardID, cellID, {
        name,
        viewProperties,
      })
    )
  }, [dispatch, dashboardID, cellID, name, viewProperties])

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

export default () => {
  const {cellID} = useParams()

  return (
    <GetResource resources={[{type: ResourceType.Cells, id: cellID}]}>
      <EditVEO />
    </GetResource>
  )
}
