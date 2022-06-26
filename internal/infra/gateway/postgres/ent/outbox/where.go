// Code generated by ent, DO NOT EDIT.

package outbox

import (
	"entgo.io/ent/dialect/sql"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int64) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int64) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int64) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int64) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int64) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int64) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int64) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int64) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int64) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// Kind applies equality check predicate on the "kind" field. It's identical to KindEQ.
func Kind(v string) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldKind), v))
	})
}

// Payload applies equality check predicate on the "payload" field. It's identical to PayloadEQ.
func Payload(v []byte) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPayload), v))
	})
}

// Consumed applies equality check predicate on the "consumed" field. It's identical to ConsumedEQ.
func Consumed(v bool) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldConsumed), v))
	})
}

// KindEQ applies the EQ predicate on the "kind" field.
func KindEQ(v string) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldKind), v))
	})
}

// KindNEQ applies the NEQ predicate on the "kind" field.
func KindNEQ(v string) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldKind), v))
	})
}

// KindIn applies the In predicate on the "kind" field.
func KindIn(vs ...string) predicate.Outbox {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Outbox(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldKind), v...))
	})
}

// KindNotIn applies the NotIn predicate on the "kind" field.
func KindNotIn(vs ...string) predicate.Outbox {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Outbox(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldKind), v...))
	})
}

// KindGT applies the GT predicate on the "kind" field.
func KindGT(v string) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldKind), v))
	})
}

// KindGTE applies the GTE predicate on the "kind" field.
func KindGTE(v string) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldKind), v))
	})
}

// KindLT applies the LT predicate on the "kind" field.
func KindLT(v string) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldKind), v))
	})
}

// KindLTE applies the LTE predicate on the "kind" field.
func KindLTE(v string) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldKind), v))
	})
}

// KindContains applies the Contains predicate on the "kind" field.
func KindContains(v string) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldKind), v))
	})
}

// KindHasPrefix applies the HasPrefix predicate on the "kind" field.
func KindHasPrefix(v string) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldKind), v))
	})
}

// KindHasSuffix applies the HasSuffix predicate on the "kind" field.
func KindHasSuffix(v string) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldKind), v))
	})
}

// KindEqualFold applies the EqualFold predicate on the "kind" field.
func KindEqualFold(v string) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldKind), v))
	})
}

// KindContainsFold applies the ContainsFold predicate on the "kind" field.
func KindContainsFold(v string) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldKind), v))
	})
}

// PayloadEQ applies the EQ predicate on the "payload" field.
func PayloadEQ(v []byte) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPayload), v))
	})
}

// PayloadNEQ applies the NEQ predicate on the "payload" field.
func PayloadNEQ(v []byte) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldPayload), v))
	})
}

// PayloadIn applies the In predicate on the "payload" field.
func PayloadIn(vs ...[]byte) predicate.Outbox {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Outbox(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldPayload), v...))
	})
}

// PayloadNotIn applies the NotIn predicate on the "payload" field.
func PayloadNotIn(vs ...[]byte) predicate.Outbox {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Outbox(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldPayload), v...))
	})
}

// PayloadGT applies the GT predicate on the "payload" field.
func PayloadGT(v []byte) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldPayload), v))
	})
}

// PayloadGTE applies the GTE predicate on the "payload" field.
func PayloadGTE(v []byte) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldPayload), v))
	})
}

// PayloadLT applies the LT predicate on the "payload" field.
func PayloadLT(v []byte) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldPayload), v))
	})
}

// PayloadLTE applies the LTE predicate on the "payload" field.
func PayloadLTE(v []byte) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldPayload), v))
	})
}

// ConsumedEQ applies the EQ predicate on the "consumed" field.
func ConsumedEQ(v bool) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldConsumed), v))
	})
}

// ConsumedNEQ applies the NEQ predicate on the "consumed" field.
func ConsumedNEQ(v bool) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldConsumed), v))
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Outbox) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Outbox) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Outbox) predicate.Outbox {
	return predicate.Outbox(func(s *sql.Selector) {
		p(s.Not())
	})
}