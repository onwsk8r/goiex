package market_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/onwsk8r/goiex/pkg/core/market"
	"github.com/onwsk8r/goiex/pkg/core/stock/fundamental"
	"github.com/onwsk8r/goiex/test/helper"
)

var _ = Describe("UpcomingSplit", func() { // nolint: dupl
	var expected []UpcomingSplit
	BeforeEach(func() {
		expected = GoldenUpcomingSplits()
	})

	It("should parse upcoming splits correctly", func() {
		var res []UpcomingSplit
		helper.TestdataFromJSON("core/market/upcoming_splits.json", &res)
		Expect(res).To(ConsistOf(expected))
	})

	Describe("Validate()", func() {
		It("should succeed if the UpcomingSplit is valid", func() {
			for idx := range expected {
				Expect(expected[idx].Validate()).To(Succeed())
			}
		})
		It("should return an error if the Symbol is empty", func() {
			expected[0].Symbol = ""
			Expect(expected[0].Validate()).To(MatchError("symbol is missing"))
		})
		It("should return an error if the DeclaredDate is zero valued", func() {
			expected[0].DeclaredDate = time.Time{}
			Expect(expected[0].Validate()).To(MatchError("declared date is missing"))
		})
		It("should return an error if the ExDate is zero valued", func() {
			expected[0].ExDate = time.Time{}
			Expect(expected[0].Validate()).To(MatchError("ex date is missing"))
		})
	})
})

var _ = XDescribe("UpcomingSplit Golden", func() {
	It("should load the golden file", func() {
		loc, err := time.LoadLocation("UTC")
		Expect(err).ToNot(HaveOccurred())
		golden := []UpcomingSplit{
			UpcomingSplit{
				Symbol: "MBCN",
				Split: fundamental.Split{
					ExDate:       time.Date(2019, time.November, 18, 0, 0, 0, 0, loc),
					DeclaredDate: time.Date(2019, time.October, 13, 0, 0, 0, 0, loc),
					Ratio:        0.5,
					ToFactor:     2,
					FromFactor:   1,
					Description:  "l-S i-op2r1tf",
				},
			},
			UpcomingSplit{
				Symbol: "CVLY",
				Split: fundamental.Split{
					ExDate:       time.Date(2019, time.December, 19, 0, 0, 0, 0, loc),
					DeclaredDate: time.Date(2019, time.October, 18, 0, 0, 0, 0, loc),
					Ratio:        0.998591,
					ToFactor:     22,
					FromFactor:   20,
					Description:  "il pr0f2S--1t2o",
				},
			},
			UpcomingSplit{
				Symbol: "CRPYF",
				Split: fundamental.Split{
					ExDate:       time.Date(2020, time.January, 21, 0, 0, 0, 0, loc),
					DeclaredDate: time.Date(2020, time.January, 29, 0, 0, 0, 0, loc),
					Ratio:        10,
					ToFactor:     1,
					FromFactor:   10,
					Description:  "p10t1o-firls -",
				},
			},
		}
		helper.ToGolden("upcoming_splits", golden)
	})
})
