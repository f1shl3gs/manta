const Username = 'admin'
const Password = 'password'
const Organization = 'test'

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

export const setupUser = (): Cypress.Chainable => {
  return cy
    .request({
      method: 'POST',
      url: `/api/v1/setup`,
      body: {
        username: Username,
        password: Password,
        organization: Organization,
      },
    })
    .then(resp => {
      expect(resp.status).eq(201)
    })
}

export const signin = (): Cypress.Chainable => {
  return cy
    .request({
      method: 'POST',
      url: `/api/v1/signin`,
      body: {
        username: Username,
        password: Password,
      },
    })
    .then(resp => {
      expect(resp.status).eq(200)
      return resp
    })
}

/* eslint-disable */
// general
Cypress.Commands.add('flush', flush)

// Account
Cypress.Commands.add('setupUser', setupUser)
Cypress.Commands.add('signin', signin)
