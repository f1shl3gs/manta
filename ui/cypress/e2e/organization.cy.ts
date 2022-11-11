describe('Organization', () => {
  beforeEach(() => {
    cy.setup().then(() => cy.visit('/'))
  })

  it('Create and switch', () => {
    // create one dashboard
    cy.getByTestID('nav-item-dashboard').click()
    cy.getByTestID('button-create-dashboard').first().click()
    cy.getByTestID('nav-item-dashboard').click()
    cy.getByTestID('dashboard-card').should('have.length', 1)

    // create another organization
    cy.getByTestID('tree-nav-user').click()
    cy.getByTestID('create-org').click()

    cy.getByTestID('create-org-name-input').type('other')
    cy.getByTestID('create-org-form-create').click() // after create we shall redirect to new org

    // no dashboard shall be found
    cy.getByTestID('nav-item-dashboard').click()
    cy.getByTestID('dashboard-card').should('have.length', 0)
  })
})
