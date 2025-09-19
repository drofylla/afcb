package main

import (
	"errors"
	"fmt"

	"github.com/sqids/sqids-go"
)

type Contact struct {
	ID          string
	ContactType string
	FirstName   string
	LastName    string
	Email       string
	Phone       string
}

type Contacts []Contact

var sqid *sqids.Sqids
var idCounter = 1

func init() {
	var err error
	sqid, err = sqids.New(sqids.Options{
		MinLength: 4,
	})
	if err != nil {
		panic(err)
	}
}

func newID() string {
	nums := []uint64{uint64(idCounter)}
	idCounter++
	id, _ := sqid.Encode(nums)
	return id
}

func (contacts *Contacts) New(contactType, firstName, lastName, email, phone string) {
	contact := Contact{
		ID:          newID(),
		ContactType: contactType,
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		Phone:       phone,
	}
	*contacts = append(*contacts, contact)
}

func (contacts *Contacts) ValidateID(id string) error {
	for _, contact := range *contacts {
		if contact.ID == id {
			return nil
		}
	}
	err := errors.New("invalid UUID: not found in contacts")
	fmt.Println(err)
	return err
}

func (contacts *Contacts) Delete(id string) error {
	if err := contacts.ValidateID(id); err != nil {
		return err
	}

	for i, c := range *contacts {
		if c.ID == id {
			*contacts = append((*contacts)[:i], (*contacts)[i+1:]...)
			return nil
		}
	}
	return errors.New("ID not found, unable to delete")
}
