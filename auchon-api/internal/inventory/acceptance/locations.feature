Feature: Organise places and the locations inside them
  To find things at home, I organise them into places, like the basement,
  and the areas and containers inside them, like a shelf and the boxes on it.

  Scenario: Add the basement as a place
    When I add a place called "Basement"
    Then my places are "Basement"

  Scenario: Places can grow beyond the basement
    Given I have a place called "Basement"
    When I add a place called "Garage"
    Then my places are "Basement, Garage"

  Scenario: A place needs a name
    When I add a place called ""
    Then I am told the name is missing
    And I have no places

  Scenario: Put a shelf in the basement and a box on the shelf
    Given I have a place called "Basement"
    When I add an area called "Metal shelf" in "Basement"
    And I add a container called "Cables" in "Metal shelf"
    Then "Basement" holds "Metal shelf"
    And "Metal shelf" holds "Cables"

  Scenario: A box needs somewhere to be
    When I add a container called "Cables" in "Attic"
    Then I am told "Attic" does not exist

  Scenario: Rename a box
    Given I have a place called "Basement"
    And "Basement" holds a container called "Box 1"
    When I rename "Box 1" to "Christmas lights"
    Then "Basement" holds "Christmas lights"

  Scenario: Move a box to another shelf
    Given I have a place called "Basement"
    And "Basement" holds an area called "Shelf A"
    And "Basement" holds an area called "Shelf B"
    And "Shelf A" holds a container called "Tools"
    When I move "Tools" to "Shelf B"
    Then "Shelf B" holds "Tools"
    And "Shelf A" is empty

  Scenario: A box cannot go inside a box it contains
    Given I have a place called "Basement"
    And "Basement" holds an area called "Shelf A"
    And "Shelf A" holds a container called "Tools"
    And "Tools" holds a container called "Screws"
    When I move "Tools" to "Screws"
    Then I am told a location cannot go inside itself
    And "Shelf A" holds "Tools"

  Scenario: A place stays at the top
    Given I have a place called "Basement"
    And I have a place called "Garage"
    When I move "Garage" to "Basement"
    Then I am told a place cannot be moved
    And my places are "Basement, Garage"

  Scenario: Remove an empty box
    Given I have a place called "Basement"
    And "Basement" holds a container called "Old box"
    When I remove "Old box"
    Then "Basement" is empty

  Scenario: A box that still holds things cannot be removed
    Given I have a place called "Basement"
    And "Basement" holds a container called "Tools"
    And "Tools" holds a container called "Screws"
    When I remove "Tools"
    Then I am told "Tools" is not empty
    And "Basement" holds "Tools"
