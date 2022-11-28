import {LocalStorage} from 'src/types/localStorage'

const StateKey = 'state'

export const loadLocalStorage = (): LocalStorage => {
  try {
    const data = localStorage.getItem(StateKey) ?? '{}'
    const state = JSON.parse(data)

    return normalizeGetLocalStorage(state)
  } catch (err) {
    console.error(`Load local settings failed, ${err}`)
  }

  return {
    inPresentationMode: false
  }
}

const normalizeGetLocalStorage = (state: LocalStorage): LocalStorage => {
  let newState = state

  if (state.inPresentationMode) {
    newState = {...newState, inPresentationMode: state.inPresentationMode}
  }

  return newState
}

export const saveToLocalStorage = (state: LocalStorage): void => {
  try {
    window.localStorage.setItem(
      StateKey,
      JSON.stringify(state)
    )
  } catch (err) {
    console.error(`Unable to save state to local storage: ${err}`)
  }
}
