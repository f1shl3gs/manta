// Libraries
import React, {FunctionComponent, useCallback, useEffect} from 'react'

// Components
import ViewEditorOverlayHeader from 'src/dashboards/components/ViewEditorOverlayHeader'
import TimeMachine from 'src/timeMachine/components/TimeMachine'
import {Page} from '@influxdata/clockface'
import GetResource from 'src/resources/components/GetResource'

// Hooks
import {useDispatch, useSelector} from 'react-redux'
import {useNavigate, useParams} from 'react-router-dom'

// Types
import {ResourceType} from 'src/types/resources'
import {resetTimeMachine, setViewName} from 'src/timeMachine/actions'
import {AppState} from 'src/types/stores'

// Actions
import {setTimeMachineFromCell} from 'src/dashboards/actions/thunks'
import {updateCell} from 'src/cells/actions/thunk'

// Selectors
import {getTimeMachine} from 'src/timeMachine/selectors'
import {loadView} from 'src/timeMachine/actions/thunks'

const EditVEO: FunctionComponent = () => {
  const dispatch = useDispatch()
  const {dashboardID, cellID} = useParams()
  const navigate = useNavigate()
  const {name, viewProperties} = useSelector((state: AppState) => {
    const timeMachine = getTimeMachine(state)
    return {
      name: timeMachine.name,
      viewProperties: timeMachine.viewProperties,
    }
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

  const setName = (text: string) => {
    dispatch(setViewName(text))
  }

  useEffect(() => {
    dispatch(setTimeMachineFromCell(cellID))
    dispatch(loadView())

    return () => {
      dispatch(resetTimeMachine())
    }
  }, [dispatch, cellID])

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
