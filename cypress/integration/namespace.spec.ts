import { contains } from 'cypress/types/jquery';
import { exists } from 'fs';
import { TopLevelMenu } from '~/cypress/integration/util/toplevelmenu';
import { Epinio } from '~/cypress/integration/util/epinio';

Cypress.config();
describe('Namespace testing', () => {
  const topLevelMenu = new TopLevelMenu();
  const epinio = new Epinio();
  const app_name = new String("testapp");

  beforeEach(() => {
    cy.login();
    cy.visit('/home');
    topLevelMenu.openIfClosed();
    epinio.epinioIcon().should('exist');
    epinio.accessEpinioMenu(Cypress.env('cluster'));
    // Make sure the Epinio nav menu is correct
    epinio.checkEpinioNav()
  });
  
  it('Create namespace', () => {
    cy.contains('Namespaces').click();
    cy.get('.m-0').should('contain', 'Namespaces');
    // Adding cy.wait to make the test passed with headless mode
    // Could be improved
    //cy.wait(2000);
    cy.contains('Create', {timeout: 4000}).click();
    cy.get('.labeled-input.create').type('mynamespace');
    cy.get('.card-actions .role-primary').click();
    // Adding cy.wait to wait for namespace creation
    // Could be improved
    //cy.wait(4000);
    cy.contains('mynamespace').should('be.visible');
  });

  it('Push an app into mynamespace', () => {
    // Should be a function later
    cy.contains('Applications').click();
    cy.get('.m-0').should('contain', 'Applications');
    cy.contains('Create').click();
    cy.get('.input-string > .labeled-input').type(app_name);
    cy.contains('Next').click();
    // Upload the test app
    cy.get('input[type="file"]').attachFile({filePath: 'sample-app.tar.gz', encoding: 'base64', mimeType: 'application/octet-stream'});
    cy.contains('Next').click();
    cy.get('.controls-row').contains('Create').click();
    // Check if all steps passed
    cy.get(':nth-child(1) > .col-badge-state-formatter > .badge-state').should('contain', 'Success').should('be.visible');
    cy.get(':nth-child(2) > .col-badge-state-formatter > .badge-state').should('contain', 'Success').should('be.visible');
    cy.get(':nth-child(3) > .col-badge-state-formatter > .badge-state', {timeout:120000}).should('contain', 'Success').should('be.visible');
    cy.get(':nth-child(4) > .col-badge-state-formatter > .badge-state', {timeout:120000}).should('contain', 'Success').should('be.visible');
    cy.get('.controls-row').contains('Done').click();
    // Should be another function to check an app
    // Make sure the app is in running state
    cy.get('.primaryheader', {timeout: 5000}).should('contain', app_name).and('contain', 'Running');
    // Make sure all app instances are up
    cy.get('.numbers', {timeout: 60000}).should('contain', '100%');
    cy.contains('Namespace: mynamespace').should('be.visible');
    var prefix_url = "https://";
    var app_url = prefix_url.concat(app_name, ".", Cypress.env('system_domain'));
    cy.contains(app_url).should('be.visible');
    // Other checks can be added
  });

  it('Delete namespace', () => {
    cy.contains('Namespaces').click();
    cy.get('.m-0').should('contain', 'Namespaces');
    cy.contains('mynamespace').click();
    cy.contains('Delete').click();
    cy.get('#confirm').type('mynamespace');
    cy.get('.card-container').contains('Delete').click();
    cy.contains('mynamespace', {timeout: 60000}).should('not.exist');
    // Make sure the app is also deleted
    cy.contains('Applications').click();
    cy.contains(app_name).should('not.exist');
  });
});
