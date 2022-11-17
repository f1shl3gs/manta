import {Organization} from '../../src/types/Organization'

export const DefaultUsername = 'admin'
export const DefaultPassword = 'pass'
export const DefaultOrganization = 'test'

export const flush = (): Cypress.Chainable => {
  return cy
    .request({
      method: 'GET',
      url: '/debug/flush',
    })
    .then(resp => {
      expect(resp.status).to.eq(200)
      return resp
    })
}

export const setup = (): Cypress.Chainable => {
  return cy
    .flush()
    .request({
      method: 'POST',
      url: `/api/v1/setup`,
      body: {
        username: DefaultUsername,
        password: DefaultPassword,
        organization: DefaultOrganization,
      },
    })
    .then(resp => {
      expect(resp.status).eq(200)

      const org = resp.body.org as Organization
      return cy.wrap(org).as('org')
    })
}

export const signin = (): Cypress.Chainable => {
  return cy
    .request({
      method: 'POST',
      url: `/api/v1/signin`,
      body: {
        username: DefaultUsername,
        password: DefaultPassword,
      },
    })
    .then(resp => {
      expect(resp.status).eq(200)
      return resp
    })
}

// DOM node getters
export const getByTestID = (
  dataTest: string,
  options?: Partial<
    Cypress.Loggable & Cypress.Timeoutable & Cypress.Withinable & Cypress.Shadow
  >
): Cypress.Chainable => {
  return cy.get(`[data-testid="${dataTest}"]`, options)
}

/* eslint-disable */
// general
Cypress.Commands.add('flush', flush)
Cypress.Commands.add('getByTestID', getByTestID)

// Account
Cypress.Commands.add('setup', setup)
Cypress.Commands.add('signin', signin)
