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
			<button hx-delete="/contacts/{{.ID}}"
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
