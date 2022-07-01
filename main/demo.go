package main

import (
	_ "github.com/go-kivik/couchdb/v3"
	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID    uuid.UUID
	Realm string
}

func main() {
	currentUser := User{
		ID:    uuid.NewV4(),
		Realm: "skytala",
	}
	p := Person{
		Name: "Tom22",
	}
	store := NewProtoStore("http://admin:admin@localhost:5984/")
	store.Store(&currentUser, &p)
	store.All(&currentUser, &Person{})
}
