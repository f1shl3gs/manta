import React, {useCallback} from 'react'
import {useHistory, useParams} from 'react-router-dom'
import useFetch, {CachePolicies} from 'shared/useFetch'

import {SpinnerContainer, TechnoSpinner} from '@influxdata/clockface'

import {useOtcl, useOtcls} from 'otcls/state'
import OtclForm from './OtclForm'
import {OtclOverlay} from './OtclOverlay'
import remoteDataState from '../../utils/rds'

const useEditor = () => {
  const {otclID} = useParams<{otclID: string}>()
  const {reload} = useOtcls()
  const history = useHistory()
  const {otcl, setOtcl} = useOtcl()
  const {data, loading, error, patch} = useFetch(
    `/api/v1/otcls/${otclID}`,
    {
      cachePolicy: CachePolicies.NO_CACHE,
      interceptors: {
        response: async ({response}) => {
          setOtcl(response.data)
          return response
        },
      },
    },
    [otclID]
  )

  return {
    submit: () => {
      return patch(otcl).then(() => {
        reload()
        history.goBack()
      })
    },
    rds: remoteDataState(data, error, loading),
  }
}

const OtclEdit: React.FC = () => {
  const {submit, rds} = useEditor()

  const history = useHistory()
  const onDismiss = useCallback(() => history.goBack(), [])

  return (
    <OtclOverlay title="Edit Otcl Config" onDismiss={onDismiss}>
      <SpinnerContainer loading={rds} spinnerComponent={<TechnoSpinner />}>
        <OtclForm onSubmit={submit} onDismiss={onDismiss} />
      </SpinnerContainer>
    </OtclOverlay>
  )
}

export default OtclEdit
