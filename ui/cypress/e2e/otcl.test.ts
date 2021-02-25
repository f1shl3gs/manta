describe('OTCL', function () {
  beforeEach(() => {
    cy.flush()

    cy.signin().then(body => {
      const {
        org: {id},
      } = body

      cy.visit(`/orgs/${id}/otcls`)
    })
  })

  describe('create otcl', () => {
    it('create', () => {
      cy.getByTestID('otcl--create-button').click()

      cy.getByTestID('editor--name-input').type('name')
      cy.getByTestID('editor--desc-input').type('desc')

      cy.getByTestID('create-otcl--submit').click()

      cy.getByTestID('resource-card').should('have.length', 1)
    })
  })
})
