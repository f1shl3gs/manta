import {OnboardResult} from './support/commands'

declare global {
  namespace Cypress {
    interface Chainable {
      flush(): Chainable<Element>

      onboard(): Chainable<Response>

      signin(): Chainable<OnboardResult>

      getByTestID(id: string): Chainable<Element>

      getByInputName(name: string): Chainable<Element>
    }
  }
}
