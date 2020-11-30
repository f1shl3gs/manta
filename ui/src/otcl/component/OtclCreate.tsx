import React, {useCallback} from 'react';
import {useHistory} from 'react-router';

import OtclForm from './OtclForm';
import {OtclOverlay} from './OtclOverlay';
import {emptyOtcl, useOtcl, useOtcls} from "../state";
import useFetch from "use-http";
import {useOrgID} from "../../shared/state/organization/organization";

const OtclCreate: React.FC = () => {
  const history = useHistory();
  const orgID = useOrgID();
  const {otcl, setOtcl} = useOtcl();
  const {reload} = useOtcls();
  const {post} = useFetch(`/api/v1/otcls`, {})

  const onDismiss = useCallback(() => {
    history.goBack()
  }, []);

  const onSubmit = useCallback(() => {
    post({
      orgID: orgID,
      name: otcl.name,
      desc: otcl.desc,
      content: otcl.content,
    })
      .then(() => {
        history.goBack()
        reload()
        setOtcl(emptyOtcl)
      })

    // todo: handle error

  }, [orgID, otcl])

  return (
    <OtclOverlay title={'Create new Otcl'} onDismiss={onDismiss}>
      <OtclForm onDismiss={onDismiss} onSubmit={onSubmit}/>
    </OtclOverlay>
  );
};

export default OtclCreate;
