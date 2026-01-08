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

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles("templates/home.gohtml")
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

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles("templates/contact.gohtml")

	if err != nil {
		log.Printf("Parsing template: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "<h1>There was an internal server error</h1>")
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("Parsing template error: %v", err)
		http.Error(w, "<p>There was an error parsing the template.</p>", http.StatusInternalServerError)
		return
	}
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
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	faqContent := `
	<h1>FAQ Page</h1>

	<section>
		<article>
			<p><span style="font-weight:bold;">Q:</span> Is there a free version?</p>
			<p><span style="font-weight:bold;">A:</span> Yes! We offer a free trial for 30 days on any paid plans.</p>
		</article>

		<article>
			<p><span style="font-weight:bold;">Q:</span> What are you support hours?</p>
			<p><span style="font-weight:bold;">A:</span> We have support staff answering email 24/7, though response times may be a bit slower on weekends.</p>
		</article>

		<article>
			<p><span style="font-weight:bold;">Q:</span> How do I contact support?</p>
	<p><span style="font-weight:bold;">A:</span> Email us - <a href="mailto:support@lenslocked.com">support@lenslocked.com</a></p>
		</article>
	</section>
	`
	fmt.Fprint(w, faqContent)

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
