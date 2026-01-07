package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Welcome to this site!</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `
<h1>Contact Page</h1>
		<p>To get in touch, email me at <a href="mailto:qayyax@gmail.com">qayyax@gmail.com</a>
	</p>
		`)
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

type Router struct{}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandler(w, r)
	case "/contact":
		contactHandler(w, r)
	case "/faq":
		faqHandler(w, r)
	default:
		notFoundHandler(w, r)
	}
}

func main() {
	r := chi.NewRouter()
	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.NotFound(notFoundHandler)
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
