package main

import (
	"log"

	_ "github.com/go-kivik/couchdb/v3"
	uuid "github.com/satori/go.uuid"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

type User struct {
	ID    uuid.UUID
	Realm string
}

func person() protoreflect.ProtoMessage {
	return &Person{}
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
	persons := store.All(&currentUser, person)
	log.Println(len(persons))
	for _, person := range persons {
		log.Println(person)
	}
}
