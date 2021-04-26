import './wdyr'

import React from 'react'
import ReactDOM from 'react-dom'
import {BrowserRouter, Route, Switch} from 'react-router-dom'
import App from './App'
import reportWebVitals from './reportWebVitals'

// Styles
import 'style/kanis.scss'
import '@influxdata/clockface/dist/index.css'

import {PresentationModeProvider} from './shared/usePresentationMode'
import Authentication from './components/Authentication'
import NotFound from './components/NotFound'
import {AuthenticationProvider} from './shared/useAuthentication'
import {OrgsProvider} from './shared/useOrgs'
import combineProviders from './utils/combine'

import {Provider as FetchProvider} from 'shared/useFetch'
import {NotificationProvider} from './shared/notification/useNotification'
import {TimeRangeProvider} from './shared/useTimeRange'
import {AutoRefreshProvider} from './shared/useAutoRefresh'
import {SearchParamsProvider} from './shared/useSearchParams'

const CombinedProvider = combineProviders([
  AuthenticationProvider,
  PresentationModeProvider,
  [FetchProvider, {}],
])

declare global {
  interface Window {
    __REACT_DEVTOOLS_GLOBAL_HOOK__: any
  }
}

export function disableReactDevTools() {
  // Check if the React Developer Tools global hook exists
  if (typeof window.__REACT_DEVTOOLS_GLOBAL_HOOK__ !== 'object') {
    return
  }

  for (const prop in window.__REACT_DEVTOOLS_GLOBAL_HOOK__) {
    if (prop === 'renderers') {
      // this line will remove that one console error

      window.__REACT_DEVTOOLS_GLOBAL_HOOK__[prop] = new Map()
    } else {
      // Replace all of its properties with a no-op function or a null value
      // depending on their types

      window.__REACT_DEVTOOLS_GLOBAL_HOOK__[prop] =
        typeof window.__REACT_DEVTOOLS_GLOBAL_HOOK__[prop] === 'function'
          ? () => {}
          : null
    }
  }
}

if (process.env.NODE_ENV === 'production') {
  disableReactDevTools()
}

ReactDOM.render(
  /*<React.StrictMode>

  </React.StrictMode>,*/

  <BrowserRouter>
    <SearchParamsProvider>
      <NotificationProvider>
        <AuthenticationProvider>
          <Authentication>
            <PresentationModeProvider>
              <OrgsProvider>
                <TimeRangeProvider>
                  <AutoRefreshProvider>
                    <Switch>
                      <Route path="/" component={App} />
                      <Route component={NotFound} />
                    </Switch>
                  </AutoRefreshProvider>
                </TimeRangeProvider>
              </OrgsProvider>
            </PresentationModeProvider>
          </Authentication>
        </AuthenticationProvider>
      </NotificationProvider>
    </SearchParamsProvider>
  </BrowserRouter>,
  document.getElementById('root')
)

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals()
