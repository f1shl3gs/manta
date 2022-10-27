describe('Onboard redirect', () => {
  beforeEach(() => {
    cy.flush().then(() => cy.visit('/'))
  })

  it('should redirect', () => {
    cy.location('pathname').should('include', 'setup')
  })
})
