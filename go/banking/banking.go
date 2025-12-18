package main

import (
	"fmt"
	"time"
)

type TransactionType string

const (
	TxDeposit     TransactionType = "DEPOSIT"
	TxWithdraw    TransactionType = "WITHDRAW"
	TxTransferIn  TransactionType = "TRANSFER_IN"
	TxTransferOut TransactionType = "TRANSFER_OUT"
)

type Transaction struct {
	ID        string
	Type      TransactionType
	Amount    int
	Balance   int // Balance after transaction
	Timestamp time.Time
	Note      string // e.g., "Transfer from ACC001"
}

type Account struct {
	ID           string
	Name         string
	Balance      int
	Transactions []Transaction
	CreatedAt    time.Time
}

type Response struct {
	Success bool
	Message string
	Data    interface{} // Can hold Account, Transaction, []Transaction, etc.
}

type BankingSystem struct {
	transaction    map[string][]*Transaction // User -> []Transaction
	account        map[string]*Account
	accountSeq     int
	transactionSeq int
}

func NewBankingSystem() *BankingSystem {
	return &BankingSystem{
		transaction: make(map[string][]*Transaction),
		account:     make(map[string]*Account),
	}
}

func (b *BankingSystem) CreateAccount(name string, initialDeposit int) Response {
	if initialDeposit < 0 {
		return Response{
			Success: false,
			Message: "initial deposit must be positive",
		}
	}

	b.accountSeq++
	accountID := fmt.Sprintf("ACC%03d", b.accountSeq)
	if _, exists := b.account[accountID]; exists {
		return Response{
			Success: false,
			Message: "account already exists",
		}
	}

	newAccount := &Account{
		ID:           accountID,
		Name:         name,
		Balance:      initialDeposit,
		Transactions: make([]Transaction, 0),
		CreatedAt:    time.Now(),
	}
	b.account[accountID] = newAccount

	if initialDeposit > 0 {
		b.transactionSeq++
		transactionID := fmt.Sprintf("TX%03d", b.transactionSeq)
		transaction := &Transaction{
			ID:        transactionID,
			Type:      TxDeposit,
			Amount:    initialDeposit,
			Balance:   initialDeposit,
			Timestamp: time.Now(),
			Note:      fmt.Sprintf("Initial Deposit By %s", newAccount.ID),
		}
		b.transaction[accountID] = append(b.transaction[accountID], transaction)
	}

	return Response{
		Success: true,
		Message: "Account created",
		Data:    newAccount,
	}
}

func (b *BankingSystem) Deposit(accountID string, deposit int) Response {
	if deposit <= 0 {
		return Response{
			Success: false,
			Message: "deposit amount must be positive",
		}
	}

	if _, exists := b.account[accountID]; !exists {
		return Response{
			Success: false,
			Message: "account not exists",
		}
	}

	account := b.account[accountID]
	account.Balance += deposit

	b.transactionSeq++
	transactionID := fmt.Sprintf("TX%03d", b.transactionSeq)
	transaction := &Transaction{
		ID:        transactionID,
		Type:      TxDeposit,
		Amount:    deposit,
		Balance:   account.Balance,
		Timestamp: time.Now(),
		Note:      fmt.Sprintf("Deposit By %s", account.ID),
	}
	b.transaction[accountID] = append(b.transaction[accountID], transaction)

	return Response{
		Success: true,
		Message: "Deposit successful",
		Data:    transaction,
	}
}

func (b *BankingSystem) Withdraw(accountID string, withdraw int) Response {
	if withdraw <= 0 {
		return Response{
			Success: false,
			Message: "withdraw amount must be positive",
		}
	}

	if _, exists := b.account[accountID]; !exists {
		return Response{
			Success: false,
			Message: "account not exists",
		}
	}

	account := b.account[accountID]
	if account.Balance-withdraw < 0 {
		return Response{
			Success: false,
			Message: "insufficient balance",
		}
	}

	account.Balance -= withdraw

	b.transactionSeq++
	transactionID := fmt.Sprintf("TX%03d", b.transactionSeq)
	transaction := &Transaction{
		ID:        transactionID,
		Type:      TxWithdraw,
		Amount:    withdraw,
		Balance:   account.Balance,
		Timestamp: time.Now(),
		Note:      fmt.Sprintf("Withdraw By %s", account.ID),
	}
	b.transaction[accountID] = append(b.transaction[accountID], transaction)

	return Response{
		Success: true,
		Message: "Withdrawal successful",
		Data:    transaction,
	}
}

func (b *BankingSystem) Transfer(accountID string, toAccountID string, amount int) Response {
	account, exists := b.account[accountID]
	if !exists {
		return Response{
			Success: false,
			Message: "account not exists",
		}
	}

	toAccount, exists := b.account[toAccountID]
	if !exists {
		return Response{
			Success: false,
			Message: "target account not exists",
		}
	}

	if accountID == toAccountID {
		return Response{Success: false, Message: "cannot transfer to same account"}
	}

	if account.Balance < amount {
		return Response{Success: false, Message: "insufficient funds"}
	}

	account.Balance -= amount
	toAccount.Balance += amount

	b.transactionSeq++
	outTx := &Transaction{
		ID:        fmt.Sprintf("TX%03d", b.transactionSeq),
		Type:      TxTransferOut,
		Amount:    amount,
		Balance:   account.Balance,
		Timestamp: time.Now(),
		Note:      fmt.Sprintf("Transfer to %s", toAccountID),
	}
	b.transaction[accountID] = append(b.transaction[accountID], outTx)

	b.transactionSeq++
	inTx := &Transaction{
		ID:        fmt.Sprintf("TX%03d", b.transactionSeq),
		Type:      TxTransferIn,
		Amount:    amount,
		Balance:   toAccount.Balance,
		Timestamp: time.Now(),
		Note:      fmt.Sprintf("Transfer from %s", accountID),
	}
	b.transaction[toAccountID] = append(b.transaction[toAccountID], inTx)

	return Response{
		Success: true,
		Message: "Transfer successful",
		Data:    outTx,
	}
}

func (b *BankingSystem) GetBalance(accountID string) Response {
	account, exists := b.account[accountID]
	if !exists {
		return Response{
			Success: false,
			Message: "account not exists",
		}
	}

	message := fmt.Sprintf("Current balance: %d", account.Balance)
	return Response{
		Success: true,
		Message: message,
		Data:    account.Balance,
	}
}

func (b *BankingSystem) GetTransactionHistory(accountID string) Response {
	history, exists := b.transaction[accountID]
	if !exists {
		return Response{
			Success: false,
			Message: "transaction not exists",
		}
	}

	message := fmt.Sprintf("Found %d transactions", len(history))
	return Response{Success: true, Message: message, Data: history}
}
