import {LocalStorage} from 'src/types/localStorage'
import {get} from 'lodash'

const StateKey = 'state'

export const loadLocalStorage = (): LocalStorage => {
  try {
    const data = localStorage.getItem(StateKey) || '{}'
    const state = JSON.parse(data)

    return normalizeGetLocalStorage(state)
  } catch (err) {
    console.error(`Load local settings failed, ${err}`)
  }
}

const normalizeGetLocalStorage = (state: LocalStorage): LocalStorage => {
  let newState = state

  const persisted = get(newState, 'app.persisted', false)
  if (persisted) {
    newState = {
      ...newState,
      app: normalizeApp(newState.app),
    }
  }

  return newState
}

const normalizeApp = (app: LocalStorage['app']) => {
  return {
    ...app,
    persisted: {
      ...app.persisted,
    },
  }
}

export const saveToLocalStorage = (state: LocalStorage): void => {
  try {
    window.localStorage.setItem(StateKey, JSON.stringify(state))
  } catch (err) {
    console.error(`Unable to save state to local storage: ${err}`)
  }
}
