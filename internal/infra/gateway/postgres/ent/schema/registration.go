package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Registration holds the schema definition for the Registration entity.
type Registration struct {
	ent.Schema
}

// Fields of the Registration.
func (Registration) Fields() []ent.Field {
	return []ent.Field{
		field.String("id"),
		field.String("email"),
		field.Bool("verified"),
	}
}

// Edges of the Registration.
func (Registration) Edges() []ent.Edge {
	return nil
}
