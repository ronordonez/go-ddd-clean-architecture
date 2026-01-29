package domain

import (
	"fmt"
)

type Price struct {
	value float64
}

func NewPrice(value float64) (Price, error) {
	if value <= 0 {
		return Price{}, ErrInvalidPrice
	}

	return Price{value: value}, nil
}

func (p Price) Value() float64 {
	return p.value
}

func (p Price) String() string {
	return fmt.Sprintf("%.2f", p.value)
}

func (p Price) Equals(other Price) bool {
	return p.value == other.value
}
