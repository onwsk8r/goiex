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

package rest

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/onwsk8r/goiex/pkg/core/stock"
)

type DividendsPeriod string

var (
	DividendsPeriod5y   DividendsPeriod = "5y"
	DividendsPeriod2y   DividendsPeriod = "2y"
	DividendsPeriod1y   DividendsPeriod = "1y"
	DividendsPeriodYTD  DividendsPeriod = "ytd"
	DividendsPeriod6m   DividendsPeriod = "6m"
	DividendsPeriod3m   DividendsPeriod = "3m"
	DividendsPeriod1m   DividendsPeriod = "1m"
	DividendsPeriodNext DividendsPeriod = "next"
)

type HistoricalPeriod string

var (
	HistoricalPeriodMax HistoricalPeriod = "max"
	HistoricalPeriod5y  HistoricalPeriod = "5y"
	HistoricalPeriod2y  HistoricalPeriod = "2y"
	HistoricalPeriod1y  HistoricalPeriod = "1y"
	HistoricalPeriodYTD HistoricalPeriod = "ytd"
	HistoricalPeriod6m  HistoricalPeriod = "6m"
	HistoricalPeriod3m  HistoricalPeriod = "3m"
	HistoricalPeriod1m  HistoricalPeriod = "1m"
	HistoricalPeriod5d  HistoricalPeriod = "5d"
)

type HistoricalIntradayPeriod string

var (
	HistoricalIntradayPeriodDate HistoricalIntradayPeriod = "date"
	HistoricalIntradayPeriod1mm  HistoricalIntradayPeriod = "1mm"
)

type SplitsPeriod string

var (
	SplitsPeriod5y SplitsPeriod = "5y"
)

// Stock exposes methods for calling Stock endoints.
// https://iexcloud.io/docs/api/#stocks-equities
type Stock struct {
	client *resty.Client
}

// NewStock creates a new Stock with the given client
// https://iexcloud.io/docs/api/#stocks-equities
func NewStock(client *resty.Client) *Stock {
	return &Stock{client: client}
}

// Dividends provides basic dividend data for US equities, ETFs, and Mutual Funds for the last 5 years.
// https://iexcloud.io/docs/api/#dividends-basic
func (s *Stock) Dividends(ctx context.Context, symbol string,
	period DividendsPeriod) (res []stock.Dividend, err error) {
	var params = map[string]string{"symbol": symbol, "range": string(period)}
	_, err = s.client.R().SetContext(ctx).SetPathParams(params).SetResult(&res).
		Get("/{version}/stock/{symbol}/dividends/{range}")
	return
}

// Earnings data for a given company including the actual EPS, consensus, and fiscal period.
// Available quarterly (last 4 quarters) and annually (last 4 years).
// https://iexcloud.io/docs/api/#earnings
func (s *Stock) Earnings(ctx context.Context, symbol string,
	params map[string]string) (res []stock.Earning, err error) {
	var pathParams = map[string]string{"symbol": symbol}
	_, err = s.client.R().SetContext(ctx).SetPathParams(pathParams).SetQueryParams(params).
		SetResult(&res).Get("/{version}/stock/{symbol}/earnings")
	return
}

// Historical returns adjusted and unadjusted historical data for up to 15 years.
// https://iexcloud.io/docs/api/#historical-prices
func (s *Stock) Historical(ctx context.Context, symbol string, period HistoricalPeriod,
	params map[string]string) (res []stock.Historical, err error) {
	var pathParams = map[string]string{"symbol": symbol, "range": string(period)}
	_, err = s.client.R().SetContext(ctx).SetPathParams(pathParams).SetQueryParams(params).
		SetResult(&res).Get("/{version}/stock/{symbol}/chart/{range}")
	return
}

func (s *Stock) HistoricalIntraday(ctx context.Context, symbol string, period HistoricalIntradayPeriod,
	params map[string]string) (res []stock.Intraday, err error) {
	var pathParams = map[string]string{"symbol": symbol, "range": string(period)}
	_, err = s.client.R().SetContext(ctx).SetPathParams(pathParams).SetQueryParams(params).
		SetResult(&res).Get("/{version}/stock/{symbol}/chart/{range}")
	return
}

// PreviousDay returns previous day adjusted price data for one or more stocks.
// https://iexcloud.io/docs/api/#previous-day-price
func (s *Stock) PreviousDay(ctx context.Context, symbol string) (res *stock.Historical, err error) {
	var params = map[string]string{"symbol": symbol}
	_, err = s.client.R().SetContext(ctx).SetPathParams(params).SetResult(&res).
		Get("/{version}/stock/{symbol}/previous")
	return
}

func (s *Stock) PreviousDayMarket(ctx context.Context) (res []stock.Historical, err error) {
	_, err = s.client.R().SetContext(ctx).SetResult(&res).Get("/{version}/stock/market/previous")
	return
}

func (s *Stock) Splits(ctx context.Context, symbol string,
	period SplitsPeriod) (res []stock.Split, err error) {
	var params = map[string]string{"symbol": symbol, "range": string(period)}
	_, err = s.client.R().SetContext(ctx).SetPathParams(params).SetResult(&res).
		Get("/{version}/stock/{symbol}/splits/{range}")
	return
}
