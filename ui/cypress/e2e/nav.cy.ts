describe('NavMenu', () => {
  beforeEach(() => {
    cy.setupUser()
      .then(() => cy.visit('/'))
  })

  it('User profile pop out', () => {
    cy.getByTestID('tree-nav-user').click()

    // user element should popout
    cy.contains('Members')
    cy.contains('About')
    cy.contains('Switch organization')
    cy.contains('Logout')
  })
})
