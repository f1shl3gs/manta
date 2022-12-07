import {Step} from 'src/setup/reducers'

export const SET_STEP = 'SET_STEP'
export const SET_SETUP_PARAMS = 'SET_SETUP_PARAMS'

export const setStep = (step: Step) =>
  ({
    type: SET_STEP,
    step,
  } as const)

export const setSetupParams = ({
  username,
  password,
  organization,
}: {
  username?: string
  password?: string
  organization?: string
}) =>
  ({
    type: SET_SETUP_PARAMS,
    username,
    password,
    organization,
  } as const)

export type Action =
  | ReturnType<typeof setStep>
  | ReturnType<typeof setSetupParams>
