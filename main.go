package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"text/template"

	"github.com/gorilla/mux"
)

var contacts Contacts

var emailRegex = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

const dataFile = "contacts.json"

// card template
var cardTmpl = template.Must(template.New("card").Parse(`
	<div class="card" id="contact-{{.ID}}">
		<div class="details">
			<span class="id">{{.ID}}</span><br>
			<strong class="name">{{.FirstName}} {{.LastName}}</strong><br>
			<span class="type">{{.ContactType}}</span><br>
			<span class="details">{{.Email}} | {{.Phone}}</span>
		</div>
		<div class="actions">
			<button class="edit-btn"
				onclick="openEditModal('{{.ID}}', '{{.ContactType}}', '{{.FirstName}}', '{{.LastName}}', '{{.Email}}', '{{.Phone}}')"
				title="Edit">
				<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="black" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
					<path d="M12 20h9"/>
					<path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4 12.5-12.5z"/>
				</svg>
			</button>
			<button class="delete-btn"
					hx-delete="/contacts/{{.ID}}"
					hx-target="#contact-{{.ID}}"
					hx-swap="outerHTML"
					title="Delete">
				<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="black" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
					<polyline points="3 6 5 6 21 6"/>
					<path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6m5 0V4a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2v2"/>
					<line x1="10" y1="11" x2="10" y2="17"/>
					<line x1="14" y1="11" x2="14" y2="17"/>
				</svg>
			</button>
		</div>
	</div>
`))

func renderCard(w http.ResponseWriter, c Contact) {
	w.Header().Set("Content-Type", "text/html")
	cardTmpl.Execute(w, c)
}

func getContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	for _, c := range contacts {
		cardTmpl.Execute(w, c)
	}
}

func addContact(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := r.FormValue("id")
	contactType := r.FormValue("ContactType")
	firstName := r.FormValue("FirstName")
	lastName := r.FormValue("LastName")
	email := r.FormValue("Email")
	phone := r.FormValue("Phone")

	// If ID is provided, it's an update, not a new contact
	if id != "" {
		updates := map[string]string{
			"ContactType": contactType,
			"FirstName":   firstName,
			"LastName":    lastName,
			"Email":       email,
			"Phone":       phone,
		}

		if err := contacts.UpdateContact(id, updates); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		// Find and return the updated contact
		for _, c := range contacts {
			if c.ID == id {
				if err := contacts.SaveToFile(dataFile); err != nil {
					http.Error(w, "failed to save contacts: "+err.Error(), http.StatusInternalServerError)
					return
				}
				renderCard(w, c)
				return
			}
		}

		http.Error(w, "contact not found after update", http.StatusNotFound)
		return
	}

	if !emailRegex.MatchString(email) {
		http.Error(w, "Invalid email address. Must contain @ and a valid domain.", http.StatusBadRequest)
		return
	}

	// This is a new contact
	contacts.New(contactType, firstName, lastName, email, phone)

	// ✅ persist after adding
	if err := contacts.SaveToFile(dataFile); err != nil {
		http.Error(w, "failed to save contacts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	newContact := contacts[len(contacts)-1]
	renderCard(w, newContact)
}

func updateContact(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Printf("UPDATE request received for id: %s, form values: %+v\n", id, r.Form)

	updates := map[string]string{
		"ContactType": r.FormValue("ContactType"),
		"FirstName":   r.FormValue("FirstName"),
		"LastName":    r.FormValue("LastName"),
		"Email":       r.FormValue("Email"),
		"Phone":       r.FormValue("Phone"),
	}

	fmt.Printf("Attempting to update contact %s with: %+v\n", id, updates)

	if err := contacts.UpdateContact(id, updates); err != nil {
		fmt.Println("Update error:", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := contacts.SaveToFile(dataFile); err != nil {
		http.Error(w, "failed to save contacts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Find and return the updated contact
	for _, c := range contacts {
		if c.ID == id {
			fmt.Printf("Successfully updated contact: %+v\n", c)
			renderCard(w, c)
			return
		}
	}

	http.Error(w, "contact not found after update", http.StatusNotFound)
}

func deleteContact(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	fmt.Println("DELETE request received for id:", id)

	if err := contacts.Delete(id); err != nil {
		fmt.Println("Delete error:", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := contacts.SaveToFile(dataFile); err != nil {
		http.Error(w, "failed to save contacts: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return empty content - HTMX will remove the element
	w.WriteHeader(http.StatusOK)
}

// Debug: return contacts as JSON (IDs + names) so you can inspect what's stored
func debugIDs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type simple struct {
		ID        string `json:"id"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}
	out := make([]simple, 0, len(contacts))
	for _, c := range contacts {
		out = append(out, simple{ID: c.ID, FirstName: c.FirstName, LastName: c.LastName})
	}
	json.NewEncoder(w).Encode(out)
}

func debugContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contacts)
}

// In main(), add this route:

func main() {
	if err := contacts.LoadFromFile(dataFile); err != nil {
		fmt.Println("Error loading contacts:", err)
	}

	router := mux.NewRouter()

	router.HandleFunc("/debug/contacts", debugContacts).Methods("GET")
	router.HandleFunc("/contacts", getContacts).Methods("GET")
	router.HandleFunc("/contacts", addContact).Methods("POST")
	router.HandleFunc("/contacts/{id}", updateContact).Methods("POST") // Changed from PUT to POST
	router.HandleFunc("/contacts/{id}", deleteContact).Methods("DELETE")

	// Serve static frontend
	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/").Handler(fs)

	fmt.Println("AFCB started at http://localhost:1330")
	http.ListenAndServe(":1330", router)
}
