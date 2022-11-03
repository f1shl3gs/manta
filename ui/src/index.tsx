// Libraries
import React from 'react'
import ReactDOM from 'react-dom/client'
import {BrowserRouter, Route, Routes} from 'react-router-dom'

// Components
import App from './App'
import {PresentationModeProvider} from 'shared/usePresentationMode'
import NotFound from './NotFound'

// Styles
import '@influxdata/clockface/dist/index.css'
import 'style/manta.scss'

import reportWebVitals from './reportWebVitals'
import {LoginPage} from 'signin/LoginPage'
import Setup from 'setup/SetupWizard'

// const router = createBrowserRouter(
//   createRoutesFromElements(
//     <Route>
//       <Route path='/' element={
//         <PresentationModeProvider>
//           <App/>
//         </PresentationModeProvider>
//       } />
//       <Route path={'/signin'} element={<LoginPage/>}/>
//       <Route path={'/setup'} element={<Setup/>}/>
//       <Route element={<NotFound/>}/>
//     </Route>
//   )
// )

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
    <Routes>
      <Route
        path="/*"
        element={
          <PresentationModeProvider>
            <App />
          </PresentationModeProvider>
        }
      />

      <Route path={'/signin'} element={<LoginPage />} />
      <Route path={'/setup'} element={<Setup />} />
      <Route element={<NotFound />} />
    </Routes>
  </BrowserRouter>
)

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals()
