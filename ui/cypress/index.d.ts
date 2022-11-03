import {flush, setupUser, signin, getByTestID} from './support/commands'

declare global {
  namespace Cypress {
    interface Chainable {
      flush: typeof flush
      signin: typeof signin
      setupUser: typeof setupUser
      getByTestID: typeof getByTestID
    }
  }
}
