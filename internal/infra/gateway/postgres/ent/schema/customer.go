package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Customer holds the schema definition for the Customer entity.
type Customer struct {
	ent.Schema
}

// Fields of the Customer.
func (Customer) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").MaxLen(36),
		field.Int("version"),
		field.String("first_name").MaxLen(50).Optional(),
		field.String("last_name").MaxLen(50).Optional(),
		field.String("email").MaxLen(255),
	}
}

// Edges of the Customer.
func (Customer) Edges() []ent.Edge {
	return nil
}
