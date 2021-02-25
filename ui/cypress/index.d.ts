declare namespace Cypress {
  interface Chainable {
    flush(): Chainable<Element>
    onboard(): Chainable<Response>
    signin(): Chainable<Element>

    getByTestID(id: string): Chainable<Element>
    getByInputName(name: string): Chainable<Element>
  }
}
