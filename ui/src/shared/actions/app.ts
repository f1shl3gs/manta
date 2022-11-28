export enum ActionTypes {
  EnablePresentationMode = 'EnablePresentationMode',
  DisablePresentationMode = 'DisablePresentationMode',
  ToggleNavBarState = 'ToggleNavBarState',
}

export type Action =
  | ReturnType<typeof disablePresentationMode>
  | ReturnType<typeof enablePresentationMode>
  | ReturnType<typeof toggleNavBarState>

export const enablePresentationMode = () =>
  ({
    type: ActionTypes.EnablePresentationMode,
  } as const)

export const disablePresentationMode = () =>
  ({
    type: ActionTypes.DisablePresentationMode,
  } as const)

export const toggleNavBarState = () =>
  ({
    type: ActionTypes.ToggleNavBarState,
  } as const)
