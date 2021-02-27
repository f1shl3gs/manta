Cypress.Commands.add(
  'getByTestID',
  (id: string): Cypress.Chainable => {
    return cy.get(`[data-testid="${id}"`)
  }
)

const onboard = () => {
  return cy.fixture('onboard').then(({username, password, org}) => {
    return cy.request({
      method: 'POST',
      url: '/api/v1/setup',
      body: {
        username,
        password,
        org,
      },
    })
  })
}

Cypress.Commands.add('onboarding', onboard)

Cypress.Commands.add(
  'signin',
  (): Cypress.Chainable => {
    return cy.fixture('onboard').then(({username, password}) => {
      return onboard().then(({body}) => {
        return cy
          .request({
            method: 'POST',
            url: '/api/v1/signin',
            body: {
              username,
              password,
            },
          })
          .then(() => {
            return cy.wrap(body)
          })
      })
    })
  }
)

Cypress.Commands.add(
  'getByInputName',
  (name: string): Cypress.Chainable => {
    return cy.get(`input[name=${name}]`)
  }
)

Cypress.Commands.add(
  'flush',
  (): Cypress.Chainable => {
    return cy.request({
      method: 'GET',
      url: '/kv/flush',
    })
  }
)
