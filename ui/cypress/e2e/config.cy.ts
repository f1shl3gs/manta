import {Organization} from '../../src/types/organization'

describe('Config', () => {
  beforeEach(() => {
    cy.setup().then(() => cy.visit('/'))

    cy.getByTestID('nav-item-data').should('be.visible')
    cy.getByTestID('nav-item-data--config').click({force: true})

    cy.getByTestID('config-card').should('have.length', 0)
  })

  it('Create', () => {
    // create
    cy.getByTestID('button-create-config').first().click()
    cy.url().should('include', 'data/config/new')
    cy.getByTestID('yaml-editor').type('foo: bar{enter}')
    cy.getByTestID('create-config--button').click()

    // create should be success
    cy.getByTestID('config-card').should('have.length', 1)
  })

  describe('when a config already exist', () => {
    beforeEach(() => {
      cy.get('@org').then((org: Organization) => {
        cy.request({
          method: 'POST',
          url: '/api/v1/configs',
          body: {
            name: 'foo',
            desc: 'bar',
            orgID: org.id,
          },
        })

        cy.visit(`/orgs/${org.id}/data/config`)
      })

      // create should be success
      cy.getByTestID('config-card').should('have.length', 1)
    })

    it('delete config', () => {
      cy.getByTestID('config-card-context--delete').click()
      cy.getByTestID('context_menu-delete').click()

      // the initial one should be deleted already
      cy.getByTestID('config-card').should('have.length', 0)
    })
  })
})
