package main

import (
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

var contacts Contacts

// card template
var cardTmpl = template.Must(template.New("card").Parse(`
	<div class="card" id="contact-{{.ID}}">
		<div class="details">
		<strong>{{.FirstName}} {{.LastName}}</strong><br>
		{{.ContactType}}<br>
		{{.Email}} | {{.Phone}}
		</div>
		<div class="actions">
			<button class="edit-btn"
				onclick="openEditModal('{{.ID}}', '{{.ContactType}}', '{{.FirstName}}', '{{.LastName}}', '{{.Email}}', '{{.Phone}}')">✏️</button>
			<button hx-delete="/contacts/{{.ID}}"
					hx-target="#contact-{{.ID}}"
					hx-swap="outerHTML">🗑️</button>
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
	// var c Contact
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	contactType := r.FormValue("ContactType")
	firstName := r.FormValue("FirstName")
	lastName := r.FormValue("LastName")
	email := r.FormValue("Email")
	phone := r.FormValue("Phone")

	contacts.New(contactType, firstName, lastName, email, phone)

	// w.Header().Set("Content-Type", "text/html")
	// for _, contact := range contacts {
	// 	fmt.Fprintf(w, `
	// 		<div class="card"
	// 			<div><strong>%s %s</strong><div>
	// 			<div>%s</div>
	// 			<div>%s</div>
	// 			<button hx-delete="/contacts/%s" hx-target="#contacts" hx-swap="outerHTML">Delete</button>
	// 		</div>
	// 		`, contact.FirstName, contact.LastName, contact.Email, contact.Phone, contact.ID)
	// }

	renderCard(w, contacts[len(contacts)-1])
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
	updates := map[string]string{}

	if v := r.FormValue("FirstName"); v != "" {
		updates["FirstName"] = v
	}
	if v := r.FormValue("LastName"); v != "" {
		updates["LastName"] = v
	}
	if v := r.FormValue("Email"); v != "" {
		updates["Email"] = v
	}
	if v := r.FormValue("Phone"); v != "" {
		updates["Phone"] = v
	}
	if v := r.FormValue("ContactType"); v != "" {
		updates["ContactType"] = v
	}

	if err := contacts.UpdateContact(id, updates); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	for _, c := range contacts {
		if c.ID == id {
			renderCard(w, c)
			return
		}
	}
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/contacts", getContacts).Methods("GET")
	router.HandleFunc("/contacts", addContact).Methods("POST")
	router.HandleFunc("/contacts/{id}", deleteContact).Methods("DELETE")
	router.HandleFunc("/contacts/{id}", updateContact).Methods("PUT")

	// Serve static frontend
	fs := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/").Handler(fs)
	http.ListenAndServe(":1330", router)
}
