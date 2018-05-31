package models

import "testing"

func TestUser(t *testing.T) {

	assertCorrectFieldAssignment := func(t *testing.T, expected, outcome string) {
		if expected != outcome {
			t.Errorf("Field was incorrect: got %s and expected %s.", expected, outcome)
		}
	}

	t.Run("Properly set username", func(t *testing.T) {
		username := "golang"
		user := User{username, ""}
		assertCorrectFieldAssignment(t, username, user.Username)
	})

	t.Run("Properly set email", func(t * testing.T) {
		email := "golang@golang.org"
		user := User{"", email}
		assertCorrectFieldAssignment(t, email, user.Email)
	})

}
