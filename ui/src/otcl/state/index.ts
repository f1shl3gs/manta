import {useCallback, useState} from 'react';
import constate from "constate";
import useFetch, {CachePolicies} from "shared/hooks/useFetch";

import {Otcl} from 'types';
import remoteDataState from "utils/rds";
import {useOrgID} from "shared/state/organization/organization";

const OtclPrefix = '/api/v1/otcls'

const emptyOtcl: Otcl = {
  id: '',
  name: '',
  desc: '',
  content: '',
  created: '',
  modified: ''
}

type State = {
  orgID: string
}

const [OtclProvider, useOtcls, useOtcl] = constate(
  ({orgID}: State) => {
    const [mutation, setMutation] = useState<number>(0);
    const [otcl, setOtcl] = useState<Otcl>(emptyOtcl);

    return {
      otcl,
      setOtcl,
      orgID,
      mutation,
      setMutation
    }
  },
  // useOtcls
  value => {
    const {mutation, setMutation, orgID} = value;
    const {data, error, loading} = useFetch<[Otcl]>(`/api/v1/otcls?orgID=${value.orgID}`,
      {}, [
        orgID, mutation
      ])

    console.log('data', data)

    const rds = remoteDataState(loading, error)
    return {
      rds,
      otcls: data,
      reload: useCallback(() => setMutation(mutation + 1), [mutation])
    }
  },
  // useOtcl
  value => {
    return {
      otcl: value.otcl,
      setOtcl: value.setOtcl,
    }
  }
)

const useOtclV2 = (id: string) => {
  const orgID = useOrgID();
  const {reload} = useOtcls();
  const {otcl, setOtcl} = useOtcl();
  const {error, loading, post, patch, del} = useFetch(`/api/v1/otcls/${id}`, {
    cachePolicy: CachePolicies.NO_CACHE,
    interceptors: {
      response: async ({response, options}) => {
        if (options.method !== 'GET') {
          reload()
        }

        // note: `response.data` is equivalent to `await response.json()`
        return response // returning the `response` is important
      }
    }
  })
  const rds = remoteDataState(loading, error)

  return {
    otcl,
    rds,
    post: useCallback(() => {
      return post({
        orgID,
        name: otcl.name,
        desc: otcl.desc,
        content: otcl.content,
      })
    }, [otcl]),
    patch: useCallback(() => {
      return patch(otcl)
    }, [otcl]),
    del: useCallback(() => {
      return del()
    }, [id])
  }
}

export {
  OtclProvider,
  useOtcls,
  emptyOtcl,
  useOtcl,
}