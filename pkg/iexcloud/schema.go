package iexcloud

// Validator defines a standard interface for schematic types to determine whether they are valid.
// Each type that correlates to an IEX data type should implement this interface, returning results
// that can be used to determine whether or not the data is usable. A price type, for example, may
// verify that all numeric values are positive and the close is greater than zero.
// The returned error need only describe the invalid data.
type Validator interface {
	Validate() error
}
