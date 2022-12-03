// Libraries
import React, {FunctionComponent, useCallback, useState} from 'react'
import {useNavigate, useParams} from 'react-router-dom'

// Components
import ViewEditorOverlayHeader from 'src/dashboards/ViewEditorOverlayHeader'
import TimeMachine from 'src/timeMachine'
import PageSpinner from 'src/shared/components/PageSpinner'
import { RemoteDataState } from '@influxdata/clockface'

// Hooks
import useFetch from 'src/shared/useFetch'
import {TimeMachineProvider, useTimeMachine} from 'src/timeMachine/useTimeMachine'

interface Props {
  name: string
}

const EditVEO: FunctionComponent<Props> = ({name: cellName}) => {
  const {viewProperties} = useTimeMachine()
  const [name, setName] = useState(cellName)
  const {dashboardID, cellID} = useParams()
  const navigate = useNavigate()
  const {run: updateCell} = useFetch(`/api/v1/dashboards/${dashboardID}/cells/${cellID}`, {
    method: 'PATCH',
    onSuccess: _ => {
      navigate(-1)
    }
  })
  const handleSubmit = useCallback(() => {
    updateCell({
      name,
      viewProperties
    })
  }, [updateCell, name, viewProperties])

  return (
    <div className={'veo'}>
        <ViewEditorOverlayHeader name={name} onRename={setName} onSubmit={handleSubmit} />

        <div className={'veo-contents'}>
          <TimeMachine />
        </div>
    </div>
  )
}

export default () => {
  const {cellID, dashboardID} = useParams()
  const {data, loading} = useFetch(
    `/api/v1/dashboards/${dashboardID}/cells/${cellID}`
  )

  if (loading === RemoteDataState.Loading) {
    return <></>
  }

  return (
    <PageSpinner loading={loading}>
      <TimeMachineProvider viewProperties={data.viewProperties}>
        <EditVEO name={data.name}  />
      </TimeMachineProvider>
    </PageSpinner>
  )
}
