package main

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type Contact struct {
	ID          uuid.UUID
	ContactType string
	FirstName   string
	LastName    string
	Email       string
	Phone       string
}

type Contacts []Contact

func (contacts *Contacts) New(contactType, firstName, lastName, email, phone string) {
	contact := Contact{
		ID:          uuid.New(),
		ContactType: contactType,
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		Phone:       phone,
	}
	*contacts = append(*contacts, contact)
}

func (contacts *Contacts) ValidateUUID(id uuid.UUID) error {
	for _, contact := range *contacts {
		if contact.ID == id {
			return nil
		}
	}
	err := errors.New("invalid UUID: not found in contacts")
	fmt.Println(err)
	return err
}
