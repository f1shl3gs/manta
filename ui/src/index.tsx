// Libraries
import React, {lazy, Suspense} from 'react'
import ReactDOM from 'react-dom/client'
import {BrowserRouter, Route, Routes} from 'react-router-dom'

// Components
import App from './App'
import {PresentationModeProvider} from 'shared/usePresentationMode'
import SetupWrapper from 'setup/SetupWrapper'
import PageSpinner from 'shared/components/PageSpinner'

// Styles
import '@influxdata/clockface/dist/index.css'
import 'style/manta.scss'

import reportWebVitals from './reportWebVitals'

const SignInPage = lazy(() => import('signin/LoginPage'))

const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement)
root.render(
  /*
  react-custom-scrollbars not works well with react v18.
  Similar issues https://github.com/xobotyi/react-scrollbars-custom/issues/234

  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
*/
  <BrowserRouter>
    <SetupWrapper>
      <Suspense fallback={<PageSpinner />}>
        <Routes>
          <Route path={'/signin'} element={<SignInPage />} />
          <Route
            path="/*"
            element={
              <PresentationModeProvider>
                <App />
              </PresentationModeProvider>
            }
          />
        </Routes>
      </Suspense>
    </SetupWrapper>
  </BrowserRouter>
)

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals()
