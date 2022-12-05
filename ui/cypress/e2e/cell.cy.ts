describe('Dashboard', () => {
  beforeEach(() => {
    cy.setup().then(() => cy.visit('/'))

    cy.getByTestID('nav-item-dashboard').click()

    // should be empty
    cy.getByTestID('dashboard-card').should('have.length', 0)

    // create one
    cy.getByTestID('add-resource-dropdown--button')
      .should('have.length', 2)
      .first()
      .click()
    cy.getByTestID('add-resource-dropdown--new').click()

    cy.location('pathname')
      .should('include', 'orgs')
      .should('include', 'dashboards')
    cy.getByTestID('nav-item-dashboard').click()
    cy.getByTestID('dashboard-editable-name').click()
  })

  it('Add cell', () => {
    const name = 'some name'
    const query = 'queeeeeery'

    cy.getByTestID('create-cell--button').first().click()

    // fill query
    cy.get(`[data-mode-id="promql"]`).type(query)

    // set name
    cy.getByTestID('page-title').first().click()
    .getByTestID('renamable-page-title--input').type(`${name}{enter}`)

    // submit
    cy.getByTestID('submit--button').click()

    // find cell
    cy.getByTestID(`cell--draggable ${name}`).should('have.length', 1)
  })
})
