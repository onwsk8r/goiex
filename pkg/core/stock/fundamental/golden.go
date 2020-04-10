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

package fundamental

import "github.com/onwsk8r/goiex/test/helper"

// GoldenDividend returns golden data for the Dividend type
func GoldenDividend() (d Dividend) {
	helper.FromGolden("dividend", &d)
	return
}

// GoldenEarnings returns golden data for the Earning type
func GoldenEarnings() (e []Earning) {
	helper.FromGolden("earning", &e)
	return
}

// GoldenSplit returns golden data for the Split type
func GoldenSplit() (s Split) {
	helper.FromGolden("split", &s)
	return
}
