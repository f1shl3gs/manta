import React, {useCallback} from "react";
import {useHistory, useParams} from "react-router";
import useFetch, {CachePolicies} from "shared/hooks/useFetch";

import {
  SpinnerContainer,
  TechnoSpinner
} from "@influxdata/clockface";

import {useOtcl, useOtcls} from "otcl/state";
import OtclForm from "./OtclForm";
import {OtclOverlay} from "./OtclOverlay";
import remoteDataState from "../../utils/rds";

const useEditor = (id: string) => {
  const {otclID} = useParams();
  const {reload} = useOtcls();
  const history = useHistory();
  const {otcl, setOtcl} = useOtcl();
  const {loading, error, patch} = useFetch(`/api/v1/otcls/${otclID}`, {
    cachePolicy: CachePolicies.NO_CACHE,
    interceptors: {
      response: async ({response}) => {
        setOtcl(response.data)
        return response
      }
    }
  }, [otclID])

  return {
    submit: () => {
      return patch(otcl)
        .then(() => {
          reload();
          history.goBack();
        })
    },
    rds: remoteDataState(loading, error)
  }
}

const OtclEdit: React.FC = () => {
  const {otclID} = useParams();
  const {submit, rds} = useEditor(otclID)

  const history = useHistory();
  const onDismiss = useCallback(() => history.goBack(), []);

  return (
    <OtclOverlay title={'Update Otcl Config'} onDismiss={onDismiss}>
      <SpinnerContainer loading={rds} spinnerComponent={<TechnoSpinner/>}>
        <OtclForm onSubmit={submit} onDismiss={onDismiss}/>
      </SpinnerContainer>
    </OtclOverlay>
  )
}

export default OtclEdit
