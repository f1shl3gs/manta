import {DefaultUsername, DefaultPassword, DefaultOrganization} from '../support/commands';

describe('Onboard redirect', () => {
  beforeEach(() => {
    cy.flush().then(() => cy.visit('/'))
  })

  it('should redirect', () => {
    cy.location('pathname').should('eq', '/setup')
  })
})

describe('Onboard', () => {
  beforeEach(() => {
    cy.flush().then(() => {
      cy.visit('/')
    })
  })

  it('setup', () => {
    cy.getByTestID('get-start').click()

    // fill form
    cy.getByTestID('input-username').type(DefaultUsername)
    cy.getByTestID('input-password').type(DefaultPassword)
    cy.getByTestID('input-organization').type(DefaultOrganization)
    cy.getByTestID('button-next').click()

    // after setup, user will be redirect to /orgs/[orgID]
    cy.location('pathname').should('include', '/orgs/')
  })
})