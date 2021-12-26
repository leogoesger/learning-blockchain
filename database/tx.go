package database

type Tx struct {
	From  Account `json:"from"`
	To    Account `json:"to"`
	Value uint    `json:"value"`
	Data  string  `json:"data"`
}

func NewTx(from, to Account, v uint, d string) *Tx {
	return &Tx{
		From:  from,
		To:    to,
		Value: v,
		Data:  d,
	}
}
