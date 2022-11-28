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
import persistStateEnhancer from 'src/store/persistStateEnhancer'

// Reducers
import app from 'src/shared/reducers/app'
import {autoRefreshReducer} from 'src/shared/reducers/autoRefresh'
import {timeRangeReducer} from 'src/shared/reducers/timeRange'

const composeEnhancers =
  (window as any).__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose

const rootReducer = (_history: History) => (state, action) => {
  if (action.type === 'UserLoggedOut') {
    state = undefined
  }

  return combineReducers<AppState>({
    app,
    autoRefresh: autoRefreshReducer,
    timeRange: timeRangeReducer
  })(state, action)
}

function configureStore(
  initialState: LocalStorage = loadLocalStorage()
): Store<AppState> {
  const create = composeEnhancers(
    persistStateEnhancer(),
    applyMiddleware(thunkMiddleware, resizeLayout)
  )(createStore)

  return create(rootReducer(history), initialState)
}

export const getStore = () => {
  return configureStore()
}
