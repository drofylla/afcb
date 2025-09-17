package main

import (
	"fmt"

	"github.com/google/uuid"
)

func main() {
	var contacts Contacts

	contacts.New("Friend", "Alice", "Smith", "alice@example.com", "123-456")
	contacts.New("Work", "Bob", "Jones", "bob@example.com", "987-654")

	// Test with existing UUID
	existingID := contacts[0].ID
	err := contacts.ValidateUUID(existingID)
	if err == nil {
		fmt.Println("✅ UUID found.")
	}

	// Test with random UUID
	randomID := uuid.New()
	err = contacts.ValidateUUID(randomID)
	if err != nil {
		fmt.Println("❌ UUID not found.")
	}
}
