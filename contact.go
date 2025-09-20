package main

import (
	"errors"
	"fmt"
	"strings"

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

func (contacts *Contacts) UpdateContact(id string, updates map[string]string) error {
	for i, c := range *contacts {
		if c.ID == id {
			for field, value := range updates {
				key := strings.ToLower(strings.ReplaceAll(field, " ", ""))
				switch key {
				case "contacttype":
					(*contacts)[i].ContactType = value
				case "firstname":
					(*contacts)[i].FirstName = value
				case "lastname":
					(*contacts)[i].LastName = value
				case "email":
					(*contacts)[i].Email = value
				case "phone":
					(*contacts)[i].Phone = value
				default:
					return fmt.Errorf("invalid field: %s ", field)
				}
			}
			return nil
		}
	}
	return errors.New("no ID found, unable to update contact info")
}

func (c *Contacts) SaveContact(id, contactType, firstName, lastName, email, phone string) {
	if id == "" {
		// New contact → generate ID
		contact := Contact{
			ID:          newID(),
			ContactType: contactType,
			FirstName:   firstName,
			LastName:    lastName,
			Email:       email,
			Phone:       phone,
		}
		*c = append(*c, contact)
		return
	}

	// Existing contact → update in place
	for i := range *c {
		if (*c)[i].ID == id {
			(*c)[i].ContactType = contactType
			(*c)[i].FirstName = firstName
			(*c)[i].LastName = lastName
			(*c)[i].Email = email
			(*c)[i].Phone = phone
			return
		}
	}
}
