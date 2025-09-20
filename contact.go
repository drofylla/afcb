package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
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
	for i, c := range *contacts {
		if c.ID == id {
			// Use slice tricks to remove the element
			*contacts = append((*contacts)[:i], (*contacts)[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("contact id %s not found", id)
}

func (contacts *Contacts) UpdateContact(id string, updates map[string]string) error {
	fmt.Printf("Looking for contact with ID: %s\n", id)
	fmt.Printf("Available contacts: %+v\n", *contacts)

	for i, c := range *contacts {
		if c.ID == id {
			fmt.Printf("Found contact at index %d: %+v\n", i, c)
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
					return fmt.Errorf("invalid field: %s", field)
				}
			}
			fmt.Printf("Updated contact: %+v\n", (*contacts)[i])
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

// Save contacts to JSON file
func (c *Contacts) SaveToFile(filename string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// Load contacts from JSON file
func (c *Contacts) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			*c = Contacts{} // start fresh if no file
			idCounter = 1
			return nil
		}
		return err
	}

	if err := json.Unmarshal(data, c); err != nil {
		return err
	}

	// ✅ Reset idCounter based on existing IDs
	max := 0
	for _, contact := range *c {
		if contact.ID == "" {
			continue
		}
		nums := sqid.Decode(contact.ID) // returns []uint64
		if len(nums) > 0 && int(nums[0]) > max {
			max = int(nums[0])
		}
	}
	idCounter = max + 1
	fmt.Println("restored idCounter to", idCounter)
	return nil
}
