Feature: Create Store

  As a store owner
  I should be able to create new stores

  Scenario: Creating a store called "Waldorf Courses"
    Given a valid store owner
    And no store called "Waldorf Courses" exists
    When I create the store called "Waldorf Courses"
    Then a store called "Waldorf Courses" exists
