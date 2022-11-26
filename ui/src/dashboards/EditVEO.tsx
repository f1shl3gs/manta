// Libraries
import React, {FunctionComponent} from 'react'
import {useParams} from 'react-router-dom'

// Components
import {Overlay, SpinnerContainer, TechnoSpinner} from '@influxdata/clockface'
import ViewEditorOverlayHeader from 'src/dashboards/ViewEditorOverlayHeader'
import TimeMachine from 'src/visualization/TimeMachine'

// Hooks
import useFetch from 'src/shared/useFetch'
import useEscape from 'src/shared/useEscape'
import {CellProvider} from 'src/dashboards/useCell'
import {ViewOptionProvider} from 'src/shared/useViewOption'

const EditVEO: FunctionComponent = () => {
  const {cellID, dashboardId} = useParams()
  const {data, loading} = useFetch(
    `/api/v1/dashboards/${dashboardId}/cells/${cellID}`
  )

  const onDismiss = useEscape()

  return (
    <Overlay visible={true} className={'veo-overlay'}>
      <div className={'veo'}>
        <SpinnerContainer
          loading={loading}
          spinnerComponent={<TechnoSpinner />}
        >
          <CellProvider cell={data}>
            <ViewOptionProvider>
              <ViewEditorOverlayHeader onDismiss={onDismiss} />

              <div className={'veo-contents'}>
                <TimeMachine viewProperties={data?.viewProperties} />
              </div>
            </ViewOptionProvider>
          </CellProvider>
        </SpinnerContainer>
      </div>
    </Overlay>
  )
}

export default EditVEO
