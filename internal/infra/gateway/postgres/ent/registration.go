// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres/ent/registration"
)

// Registration is the model entity for the Registration schema.
type Registration struct {
	config `json:"-"`
	// ID of the ent.
	ID string `json:"id,omitempty"`
	// Email holds the value of the "email" field.
	Email string `json:"email,omitempty"`
	// Verified holds the value of the "verified" field.
	Verified bool `json:"verified,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Registration) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case registration.FieldVerified:
			values[i] = new(sql.NullBool)
		case registration.FieldID, registration.FieldEmail:
			values[i] = new(sql.NullString)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Registration", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Registration fields.
func (r *Registration) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case registration.FieldID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value.Valid {
				r.ID = value.String
			}
		case registration.FieldEmail:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field email", values[i])
			} else if value.Valid {
				r.Email = value.String
			}
		case registration.FieldVerified:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field verified", values[i])
			} else if value.Valid {
				r.Verified = value.Bool
			}
		}
	}
	return nil
}

// Update returns a builder for updating this Registration.
// Note that you need to call Registration.Unwrap() before calling this method if this Registration
// was returned from a transaction, and the transaction was committed or rolled back.
func (r *Registration) Update() *RegistrationUpdateOne {
	return (&RegistrationClient{config: r.config}).UpdateOne(r)
}

// Unwrap unwraps the Registration entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (r *Registration) Unwrap() *Registration {
	_tx, ok := r.config.driver.(*txDriver)
	if !ok {
		panic("ent: Registration is not a transactional entity")
	}
	r.config.driver = _tx.drv
	return r
}

// String implements the fmt.Stringer.
func (r *Registration) String() string {
	var builder strings.Builder
	builder.WriteString("Registration(")
	builder.WriteString(fmt.Sprintf("id=%v, ", r.ID))
	builder.WriteString("email=")
	builder.WriteString(r.Email)
	builder.WriteString(", ")
	builder.WriteString("verified=")
	builder.WriteString(fmt.Sprintf("%v", r.Verified))
	builder.WriteByte(')')
	return builder.String()
}

// Registrations is a parsable slice of Registration.
type Registrations []*Registration

func (r Registrations) config(cfg config) {
	for _i := range r {
		r[_i].config = cfg
	}
}