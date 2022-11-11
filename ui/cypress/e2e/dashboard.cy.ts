describe('Dashboard', () => {
  beforeEach(() => {
    cy.setup()
    .then(() => cy.visit('/'))

    cy.getByTestID('nav-item-dashboard').click()

    // should be empty
    cy.getByTestID('dashboard-card').should('have.length', 0)

    // create one
    cy.getByTestID('button-create-dashboard')
      .should('have.length', 2)
      .first()
      .click()
    cy.location('pathname').should('include', 'orgs').should('include', 'dashboards')
    cy.getByTestID('nav-item-dashboard').click()
  })

  it('Delete', () => {
    // delete dashboard from list
    cy.getByTestID('dashboard-card-context--delete').click()
    cy.getByTestID('context_menu-delete').click()
    cy.getByTestID('notification-success').should('have.length', 1)

    // should be empty
    cy.getByTestID('dashboard-card').should('have.length', 0)
  })

  it('Rename', () => {
    cy.getByTestID('dashboard-editable-name--button').click()
    cy.getByTestID('dashboard-editable-name--input').type('foo{enter}')
    cy.getByTestID('dashboard-editable-name').invoke('text').should('eq', 'foo')
  })

  // TODO:
  // it('Update desc', () => {})
})
