import {flush, setup, signin, getByTestID} from './support/commands'

declare global {
  namespace Cypress {
    interface Chainable {
      flush: typeof flush
      signin: typeof signin
      setup: typeof setup
      getByTestID: typeof getByTestID
    }
  }
}
