describe('Create', () => {
  beforeEach(() => {
    cy.setup().then(() => cy.visit('/'))

    cy.getByTestID('nav-item-dashboard').click()

    // should be empty
    cy.getByTestID('dashboard-card').should('have.length', 0)

    cy.location('pathname')
      .should('include', 'orgs')
      .should('include', 'dashboards')
    cy.getByTestID('nav-item-dashboard').click()
  })

  it('can import dashboard', () => {
    const content = `{
  "id": "0a66228cdb616000",
  "created": "2022-12-07T01:47:33.101807828+08:00",
  "updated": "2022-12-07T01:57:01.526551906+08:00",
  "name": "Manta",
  "desc": "Metrics of Manta",
  "orgID": "0a659bccc2aba000"
}`

    cy.getByTestID('add-resource-dropdown--button')
      .should('have.length', 2)
      .first()
      .click()
    cy.getByTestID('add-resource-dropdown--import').click()
    cy.getByTestID('import-overlay--textarea')
      .type(content)
      .getByTestID('submit-Dashboard-button')
      .click()

    cy.getByTestID('dashboard-card').should('have.length', 1)
  })
})

describe('Update', () => {
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
  })

  it('can delete', () => {
    // delete dashboard from list
    cy.getByTestID('dashboard-card-context--delete--button').click()
    cy.getByTestID('dashboard-card-context--delete--confirm-button').click()
    cy.getByTestID('notification-success').should('have.length', 1)

    // should be empty
    cy.getByTestID('dashboard-card').should('have.length', 0)
  })

  it('can rename', () => {
    cy.getByTestID('dashboard-editable-name--button').click()
    cy.getByTestID('dashboard-editable-name--input').type('foo{enter}')
    cy.getByTestID('dashboard-editable-name').invoke('text').should('eq', 'foo')
  })

  it('can rename in detail page', () => {
    const name = 'foo'

    cy.getByTestID('dashboard-editable-name').click()
    cy.getByTestID('page-title')
      .click()
      .getByTestID('renamable-page-title--input')
      .type(`${name}{enter}`)

    cy.get('@org').then((org: Organization) => {
      cy.request({
        url: `/api/v1/dashboards?orgID=${org.id}`,
      }).then(resp => {
        expect(resp.status).eq(200)
        expect(resp.body[0].name).eq(name)
      })
    })
  })

  it('can update desc', () => {
    const desc = 'barrrrr'

    cy.get(`[class="cf-resource-description--preview untitled"]`).click()
    cy.getByTestID('input-field').type(`${desc}{enter}`)

    cy.wait(1000) // wait for patch finished, redux introduce latency!?
      .get('@org')
      .then((org: Organization) => {
        cy.request({url: `/api/v1/dashboards?orgID=${org.id}`}).then(resp => {
          expect(resp.status).eq(200)
          expect(resp.body[0].desc).eq(desc)
        })
      })
  })
})
