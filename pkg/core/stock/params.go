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

package stock

import (
	"net/url"
	"strconv"
	"time"
)

// HistoricalPriceParams enumerates the parameters for the historical prices endpoint.
// This type is used as a parameter for price.Historical()
type HistoricalPriceParams struct {
	ExactDate       time.Time
	Range           string
	Sort            string
	ChartInterval   int
	ChartLast       int
	ChartCloseOnly  bool
	ChartByDay      bool
	ChartSimplify   bool
	ChangeFromClose bool
	IncludeToday    bool
}

// Values returns the contents of this struct as a url.Values
func (h *HistoricalPriceParams) Values() url.Values {
	v := url.Values{}
	if h.ChartCloseOnly {
		v.Add("chartCloseOnly", "true")
	}
	if h.ChartByDay {
		v.Add("chartByDay", "true")
	}
	if h.ChartSimplify {
		v.Add("chartSimplify", "true")
	}
	if h.ChartInterval > 0 {
		v.Add("chartInterval", strconv.Itoa(h.ChartInterval))
	}
	if h.ChangeFromClose {
		v.Add("changeFromClose", "true")
	}
	if h.ChartLast > 0 {
		v.Add("chartLast", strconv.Itoa(h.ChartLast))
	}
	if h.Range != "" {
		v.Add("range", h.Range)
	}
	if !h.ExactDate.IsZero() {
		v.Add("exactDate", h.ExactDate.Format("20060102"))
	}
	if h.Sort != "" {
		v.Add("sort", h.Sort)
	}
	if h.IncludeToday {
		v.Add("includeToday", "true")
	}
	return v
}
