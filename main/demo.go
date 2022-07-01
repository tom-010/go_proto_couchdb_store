package main

import (
	"log"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
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
		Name: "Tom",
	}
	store(&currentUser, &p)
}

func store(user *User, message protoreflect.ProtoMessage) {
	log.Println(user.ID)
	log.Println(user.Realm)
	log.Println(message.ProtoReflect().Descriptor().FullName())
	log.Println(toJson(message))
}

func toJson(message protoreflect.ProtoMessage) string {
	encoded, err := protojson.Marshal(message)
	if err != nil {
		log.Fatalf("Could not encode proto-mesage: %v", err)
	}
	return string(encoded)
}
