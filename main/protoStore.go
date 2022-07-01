package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-kivik/kivik/v3"
	uuid "github.com/satori/go.uuid"
	"google.golang.org/protobuf/encoding/protojson"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

type ProtoStore struct {
	client *kivik.Client
	ctx    context.Context
}

func NewProtoStore(url string) ProtoStore {
	client, err := kivik.New("couch", url)
	if err != nil {
		panic(err)
	}

	return ProtoStore{
		client: client,
		ctx:    context.Background(),
	}
}

func (p *ProtoStore) Store(user *User, message protoreflect.ProtoMessage) {
	docId := uuid.NewV4().String()
	doc := toMap(message)
	doc["id"] = docId
	doc["type"] = message.ProtoReflect().Descriptor().FullName()
	doc["typeVersion"] = 1
	doc["createdBy"] = user.ID
	p.db(user.Realm).Put(p.ctx, docId, doc)
}

func (p *ProtoStore) All(user *User, table protoreflect.ProtoMessage) []protoreflect.ProtoMessage {
	tableName := table.ProtoReflect().Descriptor().FullName()

	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"type": map[string]interface{}{
				"$eq": tableName,
			},
		},
	}
	encoded, err := json.Marshal(query)

	if err != nil {
		log.Fatalf("could not encode query: %v", err)
	}
	rows, err := p.db(user.Realm).Find(p.ctx, encoded)
	if err != nil {
		log.Fatalf("Could not read table %s: %v", tableName, err)
	}

	protoReader := protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}

	for rows.Next() {
		var doc map[string]interface{}
		if err := rows.ScanDoc(&doc); err != nil {
			panic(err)
		}

		jsonEncoded, err := json.Marshal(doc)
		if err != nil {
			log.Fatalf("Could not reencode json")
		}
		m := Person{}
		err = protoReader.Unmarshal(jsonEncoded, &m)
		if err != nil {
			log.Fatalf("could not read protobuf message: %v", err)
		}
		log.Printf("%s, %s", m.Id, m.Name)
	}
	return make([]protoreflect.ProtoMessage, 0)
}

func (p *ProtoStore) db(name string) *kivik.DB {
	db := p.client.DB(p.ctx, name)
	if db.Err() != nil {
		err := p.client.CreateDB(p.ctx, name)
		if err != nil {
			log.Fatalf("Could not create database %s: %v", name, err)
		}
	}
	return db
}

func toMap(message protoreflect.ProtoMessage) map[string]interface{} {
	encoded, err := protojson.Marshal(message)
	if err != nil {
		log.Fatalf("Could not encode proto-mesage: %v", err)
	}
	var res map[string]interface{}
	json.Unmarshal(encoded, &res)
	return res
}
