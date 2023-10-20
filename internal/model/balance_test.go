package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBalance(t *testing.T) {
	tests := []struct {
		data     []*Transaction
		expected float64
	}{
		{
			data: []*Transaction{
				{TransactionId: "tqZi6QapS41zcEHy1",
					TransactionType: Sale,
					OrderId:         "c66oxMaisTwJQXjD",
					Amount:          10,
					Currency:        Eur,
					Description:     "Test transaction",
					AccountId:       "001"},
				{TransactionId: "tqZi6QapS41zcEHy2",
					TransactionType: Sale,
					OrderId:         "c66oxMaisTwJQXjD",
					Amount:          10,
					Currency:        Eur,
					Description:     "Test transaction",
					AccountId:       "001"},
			},
			expected: 20,
		},
		{
			data: []*Transaction{
				{TransactionId: "tqZi6QapS41zcEHy1",
					TransactionType: Sale,
					OrderId:         "c66oxMaisTwJQXjD",
					Amount:          10,
					Currency:        Eur,
					Description:     "Test transaction",
					AccountId:       "001"},
				{TransactionId: "tqZi6QapS41zcEHy2",
					TransactionType: Credit,
					OrderId:         "c66oxMaisTwJQXjD",
					Amount:          10,
					Currency:        Eur,
					Description:     "Test transaction",
					AccountId:       "001"},
			},
			expected: 0,
		},
		{
			data: []*Transaction{
				{TransactionId: "tqZi6QapS41zcEHy1",
					TransactionType: Sale,
					OrderId:         "c66oxMaisTwJQXjD",
					Amount:          10,
					Currency:        Eur,
					Description:     "Test transaction",
					AccountId:       "001"},
				{TransactionId: "tqZi6QapS41zcEHy2",
					TransactionType: Refund,
					OrderId:         "c66oxMaisTwJQXjD",
					Amount:          10,
					Currency:        Eur,
					Description:     "Test transaction",
					AccountId:       "001"},
			},
			expected: 0,
		},
		{
			data: []*Transaction{
				{TransactionId: "tqZi6QapS41zcEHy2",
					TransactionType: Refund,
					OrderId:         "c66oxMaisTwJQXjD",
					Amount:          10,
					Currency:        Eur,
					Description:     "Test transaction",
					AccountId:       "001"},
			},
			expected: -10,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			b := NewBalance(tt.data)
			assert.Equal(t, tt.expected, b.GetAmount())
		})
	}
}
