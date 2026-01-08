package main

import (
	"fmt"
	"html/template"
	"strings"
)

type User struct {
	Name string
	Age  int
	Meta UserMeta
	Bio  string
}

type UserMeta struct {
	Visits int
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	user := User{
		Name: "Ichigo Kurosaki",
		Age:  19,
		Meta: UserMeta{
			Visits: 4,
		},
		Bio: `<script>alert("haha, you have been h4x0r3d!")</script>`,
	}

	// w := os.Stdout
	var b strings.Builder
	err = t.Execute(&b, user)
	if err != nil {
		panic(err)
	}
	out := b.String()
	fmt.Println(out)

}
