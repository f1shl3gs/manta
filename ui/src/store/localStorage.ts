// Types
import {LocalStorage} from 'src/types/localStorage'

// Utils
import {get} from 'src/shared/utils/get'

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
