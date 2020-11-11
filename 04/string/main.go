package main

import "fmt"

type user struct {
	name  string
	email string
}

func (u user) String() string {
	return fmt.Sprintf("Name:%s, Email:%s", u.name, u.email)
}
func main() {
	u := user{
		name:  "John Doe",
		email: "johndoe@example.com",
	}
	fmt.Println(u)
}
