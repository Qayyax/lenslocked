package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Parses html template based on the filepath passed
func executeTemplate(w http.ResponseWriter, filepath string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles(filepath)
	if err != nil {
		log.Printf("Parsing template: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "<p>There was an error parsing the template.</p")
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("Parsing template: %v", err)
		http.Error(w, "<p>There was an error parsing the template.</p>", http.StatusInternalServerError)
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "templates/home.gohtml")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "templates/contact.gohtml")
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `
		<h1>Page you are looking for does not exist</h1>
		<p>Click <a href="/">here</a> to return to home</p>
		`)

}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, "templates/faq.gohtml")

}

func ContactCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contactID := chi.URLParam(r, "contactID")
		contact := fmt.Sprintf(`
			<h1>You are looking for a contact</h1>
			<p>The contact id is: %v
			`, contactID)
		ctx := context.WithValue(r.Context(), "contact", contact)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getContact(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	contact := ctx.Value("contact")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, contact)
}

func main() {
	r := chi.NewRouter()
	// r.Use(middleware.Logger)
	r.Get("/", homeHandler)
	r.Route("/contact", func(r chi.Router) {
		r.Get("/", contactHandler)

		r.Route("/{contactID}", func(r chi.Router) {
			r.Use(middleware.Logger)
			r.Use(ContactCtx)
			r.Get("/", getContact)
		})
	})
	// r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.NotFound(notFoundHandler)
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
