import {
  createStore,
  compose,
  applyMiddleware,
  combineReducers,
  Store,
} from 'redux'
import thunkMiddleware from 'redux-thunk'
import {resizeLayout} from 'src/shared/middleware/resizeLayout'
import {AppState} from 'src/types/stores'
import {LocalStorage} from 'src/types/localStorage'
import {loadLocalStorage} from 'src/store/localStorage'

// Reducers
import app from 'src/shared/reducers/app'

const composeEnhancers =
  (window as any).__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose

const rootReducer = (_history: History) => (state, action) => {
  if (action.type === 'UserLoggedOut') {
    state = undefined
  }

  return combineReducers<AppState>({
    app,
  })(state, action)
}

function configureStore(
  initialState: LocalStorage = loadLocalStorage()
): Store<AppState> {
  const create = composeEnhancers(
    applyMiddleware(thunkMiddleware, resizeLayout)
  )(createStore)

  return create(rootReducer(history), initialState)
}

export const getStore = () => {
  return configureStore()
}
