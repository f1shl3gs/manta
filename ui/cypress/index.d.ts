import {flush, setupUser, signin} from './support/commands'

declare global {
  namespace Cypress {
    interface Chainable {
      flush: typeof flush
      signin: typeof signin
      setupUser: typeof setupUser
    }
  }
}
