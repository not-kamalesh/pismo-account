package types

type OperationType int16

const (
	OTInvalid OperationType = iota
	OTPurchase
	OTPurchaseWithInstallments
	OTWithdrawal
	OTCreditVoucher
)

const (
	OTStrPurchase                 = "PURCHASE"
	OTStrPurchaseWithInstallments = "PURCHASE_WITH_INSTALLMENTS"
	OTStrWithdrawal               = "WITHDRAWAL"
	OTStrCreditVoucher            = "CREDIT_VOUCHER"
	OTStrInvalid                  = "INVALID"
)

var (
	opTypeToStr = map[OperationType]string{
		OTPurchase:                 OTStrPurchase,
		OTPurchaseWithInstallments: OTStrPurchaseWithInstallments,
		OTWithdrawal:               OTStrWithdrawal,
		OTCreditVoucher:            OTStrCreditVoucher,
		OTInvalid:                  OTStrInvalid,
	}

	strToOpType = map[string]OperationType{
		OTStrPurchase:                 OTPurchase,
		OTStrPurchaseWithInstallments: OTPurchaseWithInstallments,
		OTStrWithdrawal:               OTWithdrawal,
		OTStrCreditVoucher:            OTCreditVoucher,
		OTStrInvalid:                  OTInvalid,
	}

	opTypeToEntryType = map[OperationType]EntryType{
		OTPurchase:                 ETDebit,
		OTPurchaseWithInstallments: ETDebit,
		OTWithdrawal:               ETDebit,
		OTCreditVoucher:            ETCredit,
	}
)

func (op OperationType) String() string {
	if str, ok := opTypeToStr[op]; ok {
		return str
	}
	return OTStrInvalid
}

func ParseOperationType(s string) OperationType {
	if op, ok := strToOpType[s]; ok {
		return op
	}
	return OTInvalid
}

func (op OperationType) GetEntryType() EntryType {
	if entryType, ok := opTypeToEntryType[op]; ok {
		return entryType
	}
	return ETNone
}

func (op OperationType) IsValid() bool {
	return op.String() != OTStrInvalid
}
