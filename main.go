package main

import (
	"fmt"
)

func main() {
	var contacts Contacts

	contacts.New("Friend", "Orm", "Korn", "ok@kshhh.com", "012-3456-789")
	contacts.New("Family", "Zal", "Ahm", "za@zba.com", "133-0133-013")

	fmt.Println("Contacts List")
	for _, c := range contacts {
		fmt.Printf("ID: %s | Name: %s %s | Email: %s | Phone: %s\n", c.ID, c.FirstName, c.LastName, c.Email, c.Phone)
	}
}
