package option

import "database/sql/driver"

// side represents the side of an options trade (eg call, put)
type side string

const (
	call side = "call" // nolint:deadcode,varcheck,unused
	put  side = "put"  // nolint:deadcode,varcheck,unused
)

func (s *side) Scan(value interface{}) error {
	*s = side(value.([]byte))
	return nil
}

func (s side) Value() (driver.Value, error) {
	return string(s), nil
}

// style represents the options style (American, European, N/A)
type style string

const (
	styleAmerican style = "A" // nolint:deadcode,varcheck,unused
	styleEuropean style = "E" // nolint:deadcode,varcheck,unused
	styleNA       style = "X" // nolint:deadcode,varcheck,unused
)

func (s *style) Scan(value interface{}) error {
	*s = style(value.([]byte))
	return nil
}

func (s style) Value() (driver.Value, error) {
	return string(s), nil
}

// oType represents the type of option (equity, index)
type oType string

const (
	equity oType = "equity" // nolint:deadcode,varcheck,unused
	index  style = "index"  // nolint:deadcode,varcheck,unused
)

func (t *oType) Scan(value interface{}) error {
	*t = oType(value.([]byte))
	return nil
}

func (t oType) Value() (driver.Value, error) {
	return string(t), nil
}
