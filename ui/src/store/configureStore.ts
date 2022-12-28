import {
  createStore,
  compose,
  applyMiddleware,
  combineReducers,
  Store,
} from 'redux'
import {History} from 'history'
import thunkMiddleware from 'redux-thunk'
import {resizeLayout} from 'src/shared/middleware/resizeLayout'
import {
  createRouterMiddleware,
  createRouterReducer,
} from '@lagunovsky/redux-react-router'

import {AppState} from 'src/types/stores'
import {LocalStorage} from 'src/types/localStorage'
import {loadLocalStorage} from 'src/store/localStorage'
// import persistStateEnhancer from 'src/store/persistStateEnhancer'

// Global history
import {history} from 'src/store/history'

// Reducers
import appReducer from 'src/shared/reducers/app'
import {autoRefreshReducer} from 'src/shared/reducers/autoRefresh'
import {timeRangeReducer} from 'src/shared/reducers/timeRange'
import {dashboardsReducer} from 'src/dashboards/reducers'
import {ResourceType} from 'src/types/resources'
import {organizationsReducer} from 'src/organizations/reducers'
import {notificationsReducer} from 'src/shared/reducers/notifications'
import {cellsReducer} from 'src/cells/reducers'
import {timeMachineReducer} from 'src/timeMachine/reducers'
import {configurationsReducer} from 'src/configurations/reducers'
import {usersReducer} from 'src/members/reducers'
import {scrapesReducers} from 'src/scrapes/reducers'
import {meReducer} from 'src/me/reducers'
import {setupReducer} from 'src/setup/reducers'
import {checksReducer} from 'src/checks/reducers'
import {checkBuilderReducer} from 'src/checks/reducers/builder'
import {secretsReduer} from 'src/secrets/reducers'

const composeEnhancers =
  (window as any).__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose

const rootReducer = (history: History) => (state, action) => {
  if (action.type === 'UserLoggedOut') {
    state = undefined
  }

  return combineReducers<AppState>({
    router: createRouterReducer(history),
    app: appReducer,
    autoRefresh: autoRefreshReducer,
    checkBuilder: checkBuilderReducer,
    me: meReducer,
    notifications: notificationsReducer,
    resources: combineReducers({
      [ResourceType.Cells]: cellsReducer,
      [ResourceType.Checks]: checksReducer,
      [ResourceType.Configurations]: configurationsReducer,
      [ResourceType.Dashboards]: dashboardsReducer,
      [ResourceType.Members]: usersReducer,
      [ResourceType.Organizations]: organizationsReducer,
      [ResourceType.Secrets]: secretsReduer,
      [ResourceType.Scrapes]: scrapesReducers,
    }),
    setup: setupReducer,
    timeRange: timeRangeReducer,
    timeMachine: timeMachineReducer,
  })(state, action)
}

function configureStore(
  initialState: LocalStorage = loadLocalStorage()
): Store<AppState> {
  const routerMiddleware = createRouterMiddleware(history)
  const create = composeEnhancers(
    // persistStateEnhancer(),
    applyMiddleware(thunkMiddleware, routerMiddleware, resizeLayout)
  )(createStore)

  return create(rootReducer(history), initialState)
}

export const getStore = () => {
  return configureStore()
}
