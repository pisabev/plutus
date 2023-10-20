package model

type Balance struct {
	transactions []*Transaction
}

func NewBalance(transactions []*Transaction) *Balance {
	return &Balance{
		transactions: transactions,
	}
}

// GetAmount TODO currency conversion
func (m *Balance) GetAmount() float64 {
	var total float64
	for _, r := range m.transactions {
		if r.TransactionType == Sale {
			total += r.Amount
		} else if r.TransactionType == Refund || r.TransactionType == Credit {
			total -= r.Amount
		}
	}
	return total
}
