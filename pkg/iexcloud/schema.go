// goiex: Golang interface to IEX Cloud API
// Copyright (C) 2019 Brian Hazeltine

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package iexcloud

// Validator defines a standard interface for schematic types to determine whether they are valid.
// Each type that correlates to an IEX data type should implement this interface, returning results
// that can be used to determine whether or not the data is usable. A price type, for example, may
// verify that all numeric values are positive and the close is greater than zero.
// The returned error need only describe the invalid data.
type Validator interface {
	Validate() error
}
