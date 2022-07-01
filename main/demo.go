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
	// TODO: pass context from extern
	// TODO: do a user.bind

	currentUser := User{
		ID:    uuid.NewV4(),
		Realm: "skytala",
	}
	p := Person{
		Name: "Tom22",
	}
	store := NewProtoStore("http://admin:admin@localhost:5984/")
	store.Store(&currentUser, &p)
	persons := store.Filter(&currentUser, person, map[string]interface{}{
		"id": map[string]interface{}{
			"$eq": "029fd7a4-b99a-4c99-866a-e04833b0dcfe",
		},
		"name": map[string]interface{}{
			"$eq": "Tom22",
		},
	})

	for _, person := range persons {
		if p, ok := person.(*Person); ok {
			log.Printf("%s: %s", p.Id, p.Name)
		}
	}
	log.Println(len(persons))
}
