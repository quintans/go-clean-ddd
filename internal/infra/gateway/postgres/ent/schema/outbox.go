package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Outbox holds the schema definition for the Outbox entity.
type Outbox struct {
	ent.Schema
}

// Fields of the Outbox.
func (Outbox) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.String("kind").MaxLen(50),
		field.Bytes("payload"),
		field.Bool("consumed"),
	}
}

// Edges of the Outbox.
func (Outbox) Edges() []ent.Edge {
	return nil
}
