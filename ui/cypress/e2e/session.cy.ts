import {DefaultUsername, DefaultPassword} from '../support/commands'

describe('Session', () => {
  beforeEach(() => {
    cy.setup()
  })

  it('SignIn', () => {
    cy.clearCookies()
    cy.visit('/')

    cy.location('pathname').should('eq', '/signin')
    cy.getByTestID('username-input').type(DefaultUsername)
    cy.getByTestID('password-input').type(DefaultPassword)
    cy.getByTestID('signin-button').click()

    cy.location('pathname').should('include', '/orgs/')
  })

  it('SignOut', () => {
    cy.visit('/')

    cy.getByTestID('tree-nav-user').click()
    cy.getByTestID('user-logout').click()

    cy.location('pathname').should('eq', '/signin')
  })
})
