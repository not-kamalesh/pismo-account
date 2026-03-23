package types

type EntryType string

const (
	ETNone   = ""
	ETCredit = "CREDIT"
	ETDebit  = "DEBIT"
)

func (e EntryType) IsValid() bool {
	return e == ETCredit || e == ETDebit
}
