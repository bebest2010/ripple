package data

type TxBase struct {
	hashable
	TransactionType    TransactionType
	Flags              *TransactionFlag `json:",omitempty"`
	SourceTag          *uint32          `json:",omitempty"`
	Account            Account
	Sequence           uint32
	Fee                Value
	SigningPubKey      *PublicKey      `json:",omitempty"`
	TxnSignature       *VariableLength `json:",omitempty"`
	Memos              Memos           `json:",omitempty"`
	PreviousTxnID      *Hash256        `json:",omitempty"`
	LastLedgerSequence *uint32         `json:",omitempty"`
}

type Payment struct {
	TxBase
	Destination    Account
	Amount         Amount
	SendMax        *Amount  `json:",omitempty"`
	Paths          *PathSet `json:",omitempty"`
	DestinationTag *uint32  `json:",omitempty"`
	InvoiceID      *Hash256 `json:",omitempty"`
}

type AccountSet struct {
	TxBase
	EmailHash     *Hash128        `json:",omitempty"`
	WalletLocator *Hash256        `json:",omitempty"`
	WalletSize    *uint32         `json:",omitempty"`
	MessageKey    *VariableLength `json:",omitempty"`
	Domain        *VariableLength `json:",omitempty"`
	TransferRate  *uint32         `json:",omitempty"`
	SetFlag       *uint32         `json:",omitempty"`
	ClearFlag     *uint32         `json:",omitempty"`
}

type SetRegularKey struct {
	TxBase
	RegularKey *RegularKey `json:",omitempty"`
}

type OfferCreate struct {
	TxBase
	OfferSequence *uint32 `json:",omitempty"`
	TakerPays     Amount
	TakerGets     Amount
	Expiration    *uint32 `json:",omitempty"`
}

type OfferCancel struct {
	TxBase
	OfferSequence uint32
}

type TrustSet struct {
	TxBase
	LimitAmount *Amount `json:",omitempty"`
	QualityIn   *uint32 `json:",omitempty"`
	QualityOut  *uint32 `json:",omitempty"`
}

type SetFee struct {
	TxBase
	BaseFee           uint64
	ReferenceFeeUnits uint32
	ReserveBase       uint32
	ReserveIncrement  uint32
}

type Amendment struct {
	TxBase
	Amendment Hash256
}

func (t *TxBase) GetBase() *TxBase {
	return t
}

func (t *TxBase) GetTransactionType() TransactionType {
	return t.TransactionType
}

func (t *TxBase) GetType() string {
	return txNames[t.TransactionType]
}

func (t *TxBase) MemoSymbol() string {
	if len(t.Memos) > 0 {
		return "✐"
	}
	return " "
}

func (t *TransactionWithMetaData) GetType() string {
	return "TransactionWithMetadata"
}

func (o *OfferCreate) Ratio() *Value {
	ratio, err := o.TakerPays.Value.Ratio(*o.TakerGets.Value)
	if err != nil {
		//TODO: Is this correct behaviour?
		return &zeroNonNative
	}
	return ratio
}
