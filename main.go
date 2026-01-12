package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/qayyax/lenslock/controllers"
	"github.com/qayyax/lenslock/views"
)

type User struct {
	Name      string
	Email     string
	Age       int
	Weight    float32
	School    map[string]string
	AnimeList []string
}

// Parses html template based on the filepath passed
func executeTemplate(w http.ResponseWriter, filepath string) {

	user := User{
		Name:   "Tife O",
		Email:  "qayyax@gmail.com",
		Age:    26,
		Weight: 175,
		School: map[string]string{
			"name":     "Trent University",
			"location": "Peterborough, ON",
		},
		AnimeList: []string{"Bleach", "One Piece", "Naruto", "Dragonball", "Full Metal Alchemist", "Jujutsu Kaizen"},
	}

	viewTpl, err := views.Parse(filepath)
	if err != nil {
		log.Printf("parsing template: %v", err)
		http.Error(w, "There was an error parsing the template.", http.StatusInternalServerError)
		return
	}
	viewTpl.Execute(w, user)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `
		<h1>Page you are looking for does not exist</h1>
		<p>Click <a href="/">here</a> to return to home</p>
		`)

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

	t, _ := views.Parse("templates/home.gohtml")
	homeController := controllers.StaticController(t)
	r.Get("/", homeController)

	t, _ = views.Parse("templates/contact.gohtml")
	contactController := controllers.StaticController(t)

	r.Route("/contact", func(r chi.Router) {
		r.Get("/", contactController)

		r.Route("/{contactID}", func(r chi.Router) {
			r.Use(middleware.Logger)
			r.Use(ContactCtx)
			r.Get("/", getContact)
		})
	})
	t, _ = views.Parse("templates/faq.gohtml")
	faqController := controllers.StaticController(t)
	r.Get("/faq", faqController)
	r.NotFound(notFoundHandler)
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
