export enum ActionTypes {
  EnablePresentationMode = 'EnablePresentationMode',
  DisablePresentationMode = 'DisablePresentationMode',
}

export type Action =
  | ReturnType<typeof disablePresentationMode>
  | ReturnType<typeof enablePresentationMode>

export const enablePresentationMode = () =>
  ({
    type: ActionTypes.EnablePresentationMode,
  } as const)

export const disablePresentationMode = () =>
  ({
    type: ActionTypes.DisablePresentationMode,
  } as const)
