import './wdyr'

import React from 'react'
import ReactDOM from 'react-dom'
import {BrowserRouter, Route, Switch} from 'react-router-dom'
import App from './App'
import reportWebVitals from './reportWebVitals'

import 'style/kanis.scss'
import '@influxdata/clockface/dist/index.css'
import {PresentationModeProvider} from './shared/usePresentationMode'
import Authentication from './components/Authentication'
import Signin from './components/Signin'
import NotFound from './components/NotFound'

ReactDOM.render(
  /*<React.StrictMode>

  </React.StrictMode>,*/

  <BrowserRouter>
    <Authentication>
      <PresentationModeProvider>
        <Switch>
          <Route exact path="/signin" component={Signin} />
          <Route path="/" component={App} />
          <Route component={NotFound} />
        </Switch>
      </PresentationModeProvider>
    </Authentication>
  </BrowserRouter>,
  document.getElementById('root')
)

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals()
