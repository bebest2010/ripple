package websockets

import (
	"encoding/json"
	"github.com/donovanhide/ripple/data"
	"sync/atomic"
)

var counter uint64

type Syncer interface {
	Done()
}

type Command struct {
	Id      uint64 `json:"id"`
	Command string `json:"command"`
	Response
}

type SynchronousCommand struct {
	Command
	Ready chan bool `json:"-"`
}

func (s SynchronousCommand) Done() {
	s.Ready <- true
}

// Fields that are in every json response
type Response struct {
	Id           uint64
	Type         string
	Status       string
	Error        string
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func newCommand(command string) Command {
	return Command{
		Id:      atomic.AddUint64(&counter, 1),
		Command: command,
	}
}

func newSynchonousCommand(command string) SynchronousCommand {
	return SynchronousCommand{
		Command: newCommand(command),
		Ready:   make(chan bool),
	}
}

type LedgerCommand struct {
	Command
	LedgerIndex  interface{} `json:"ledger_index"`
	Accounts     bool        `json:"accounts"`
	Transactions bool        `json:"transactions"`
	Expand       bool        `json:"expand"`
	Result       *struct {
		Ledger struct {
			LedgerSequence  uint32                `json:"ledger_index,string"`
			Accepted        bool                  `json:"accepted"`
			CloseTime       data.RippleTime       `json:"close_time"`
			Closed          bool                  `json:"closed"`
			Hash            data.Hash256          `json:"ledger_hash"`
			PreviousLedger  data.Hash256          `json:"parent_hash"`
			TotalXRP        uint64                `json:"total_coins,string"`
			AccountHash     data.Hash256          `json:"account_hash"`
			TransactionHash data.Hash256          `json:"transaction_hash"`
			Transactions    data.TransactionSlice `json:"transactions"`
		}
	} `json:"result,omitempty"`
}

// Creates new `ledger` command to request a ledger by index
func Ledger(ledger interface{}, transactions bool) *LedgerCommand {
	return &LedgerCommand{
		Command:      newCommand("ledger"),
		LedgerIndex:  ledger,
		Transactions: transactions,
		Expand:       true,
	}
}

type TxCommand struct {
	SynchronousCommand
	Transaction data.Hash256 `json:"transaction"`
	Result      *TxResult    `json:"result,omitempty"`
}

type TxResult struct {
	data.TransactionWithMetaData
	Validated bool `json:"validated"`
}

// A shim to populate the Validated field before passing
// control on to TransactionWithMetaData.UnmarshalJSON
func (txr *TxResult) UnmarshalJSON(b []byte) error {
	var extract map[string]interface{}
	if err := json.Unmarshal(b, &extract); err != nil {
		return err
	}
	txr.Validated = extract["validated"].(bool)
	return json.Unmarshal(b, &txr.TransactionWithMetaData)
}

type SubmitCommand struct {
	SynchronousCommand
	TxBlob string        `json:"tx_blob"`
	Result *SubmitResult `json:"result,omitempty"`
}

type SubmitResult struct {
	EngineResult        data.TransactionResult `json:"engine_result"`
	EngineResultCode    int                    `json:"engine_result_code"`
	EngineResultMessage string                 `json:"engine_result_message"`
	TxBlob              string                 `json:"tx_blob"`
	Tx                  interface{}            `json:"tx_json"`
}
