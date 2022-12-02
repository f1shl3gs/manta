// Why did you render
// import './wdyr'

// Libraries
import React, {lazy, Suspense} from 'react'
import ReactDOM from 'react-dom/client'
import {BrowserRouter, Route, Routes} from 'react-router-dom'

// Components
import App from 'src/App'
import {PresentationModeProvider} from 'src/shared/usePresentationMode'
import SetupWrapper from 'src/setup/SetupWrapper'
import PageSpinner from 'src/shared/components/PageSpinner'
import NotFound from 'src/NotFound'

// Styles
import '@influxdata/clockface/dist/index.css'
import 'src/style/manta.scss'
import 'react-virtualized/styles.css'

// Utils
import reportWebVitals from 'src/reportWebVitals'

// Lazy Load
const SignInPage = lazy(() => import('src/signin/LoginPage'))

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
          <Route path="*" element={<NotFound />} />
        </Routes>
      </Suspense>
    </SetupWrapper>
  </BrowserRouter>
)

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals()
