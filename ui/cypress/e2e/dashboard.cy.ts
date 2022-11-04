describe('Dashboard', () => {
  beforeEach(() => {
    cy.setupUser()
    .then(() => cy.visit('/'))

    cy.getByTestID('nav-item-dashboard').click()
  })

  it('Creat and Delete', () => {
    // should be empty
    cy.getByTestID('dashboard-card').should('have.length', 0)

    // create and redirect to
    cy.getByTestID('button-create-dashboard')
      .should('have.length', 2)
      .first()
      .click()
    cy.location('pathname').should('include', 'orgs').should('include', 'dashboards')

    // click and redirect to
    cy.getByTestID('nav-item-dashboard').click()
    cy.getByTestID('dashboard-editable-name').click()
    cy.location('pathname').should('include', 'orgs').should('include', 'dashboards')
    
    // delete dashboard from list
    cy.getByTestID('nav-item-dashboard').click()
    cy.getByTestID('dashboard-card-context--delete').click()
    cy.getByTestID('context_menu-delete').click()
    cy.getByTestID('notification-success').should('have.length', 1)

    // should be empty
    cy.getByTestID('dashboard-card').should('have.length', 0)
  })
})
