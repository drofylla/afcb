package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

var contacts Contacts

func getContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application.json")
	json.NewEncoder(w).Encode(contacts)
}

func addContact(w http.ResponseWriter, r *http.Request) {
	var c Contact
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func deleteContact(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := contacts.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func updateContact(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var updates map[string]string
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := contacts.UpdateContact(id, updates); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder().Encode("updated")
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/contacts", getContacts).Methods("GET")
	router.HandleFunc("/contacts", addContact).Methods("POST")
	router.HandleFunc("/contacts/{id}", deleteContact).Methods("DELETE")
	router.HandleFunc("/contacts/{id}", updateContact).Methods("PUT")

	http.ListenAndServe(":1330", router)

	// contacts.New("Friend", "Orm", "Korn", "ok@kshhh.com", "012-3456-789")
	// contacts.New("Family", "Zal", "Ahm", "za@zba.com", "133-0133-013")
	// contacts.New("Work", "Tay", "Swi", "ts@tas.com", "198-9198-919")

	// fmt.Println("Contacts List")
	// for _, c := range contacts {
	// 	fmt.Printf("ID: %s | Name: %s %s | Email: %s | Phone: %s\n", c.ID, c.FirstName, c.LastName, c.Email, c.Phone)
	// }

	// deleteID := contacts[1].ID
	// fmt.Printf("\nDeleting ID %s..\n", deleteID)
	// if err := contacts.Delete(deleteID); err != nil {
	// 	fmt.Println("Error:", err)
	// }

	// fmt.Println("\nUpdated Contacts List")
	// for _, c := range contacts {
	// 	fmt.Printf("ID: %s | Name: %s %s | Email: %s | Phone: %s\n", c.ID, c.FirstName, c.LastName, c.Email, c.Phone)
	// }

	// updateID := contacts[1].ID
	// fmt.Printf("\nUpdating Taylor's name...\n")

	// err := contacts.UpdateContact(updateID, map[string]string{
	// 	"First Name": "Taylor",
	// 	"Last Name":  "Swift",
	// })
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// }

	// fmt.Println("\nUpdated Contacts List")
	// for _, c := range contacts {
	// 	fmt.Printf("ID: %s | Name: %s %s | Email: %s | Phone: %s\n", c.ID, c.FirstName, c.LastName, c.Email, c.Phone)
	// }

}
