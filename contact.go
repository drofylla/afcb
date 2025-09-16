package main

import (
	"github.com/google/uuid"
)

type Contact struct {
	contactID   uuid.UUID
	contactType string
	firstName   string
	lastName    string
	email       string
	phone       string
}

func NewContact(contactID uuid.UUID, contactType string, firstName string, lastName string, email string, phone string) *Contact {
	return &Contact{
		contactID:   contactID,
		contactType: contactType,
		firstName:   firstName,
		lastName:    lastName,
		email:       email,
		phone:       phone,
	}
}
