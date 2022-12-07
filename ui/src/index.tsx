// Why did you render
// import './wdyr'

// Libraries
import React, {lazy, Suspense} from 'react'
import ReactDOM from 'react-dom/client'
import {Route, Routes} from 'react-router-dom'

// Components
import App from 'src/App'
import {Provider} from 'react-redux'
import Setup from 'src/Setup'
import PageSpinner from 'src/shared/components/PageSpinner'
import NotFound from 'src/NotFound'

// Styles
import '@influxdata/clockface/dist/index.css'
import 'src/style/manta.scss'
import 'react-virtualized/styles.css'

// Utils
import reportWebVitals from 'src/reportWebVitals'
import {getStore} from 'src/store/configureStore'
import {ReduxRouter} from '@lagunovsky/redux-react-router'

import {history} from 'src/store/history'

// Lazy Load
const SignInPage = lazy(() => import('src/signin/LoginPage'))

const root = ReactDOM.createRoot(document.getElementById('root'))
root.render(
  /*
  react-custom-scrollbars not works well with react v18.
  Similar issues https://github.com/xobotyi/react-scrollbars-custom/issues/234

  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
*/
  <Provider store={getStore()}>
    <ReduxRouter history={history}>
      <Setup>
        <Suspense fallback={<PageSpinner />}>
          <Routes>
            <Route path={'/signin'} element={<SignInPage />} />

            <Route path="/*" element={<App />} />
            <Route path="*" element={<NotFound />} />
          </Routes>
        </Suspense>
      </Setup>
    </ReduxRouter>
  </Provider>
)

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals()
