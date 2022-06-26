package customer

import (
	"github.com/google/uuid"
	"github.com/quintans/faults"
)

type CustomerID struct {
	id uuid.UUID
}

func ParseCustomerID(s string) (CustomerID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return CustomerID{}, faults.Wrap(err)
	}
	c := CustomerID{
		id: id,
	}

	return c, nil
}

func NewCustomerID() CustomerID {
	return CustomerID{
		id: uuid.New(),
	}
}

func MustParseCustomerID(
	id string,
) CustomerID {
	c, err := ParseCustomerID(id)
	if err != nil {
		panic(err)
	}
	return c
}

func (c CustomerID) Id() uuid.UUID {
	return c.id
}

func (c CustomerID) IsZero() bool {
	return c == CustomerID{}
}

func (c CustomerID) String() string {
	return c.id.String()
}
