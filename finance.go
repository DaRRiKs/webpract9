package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

const (
	incomeType  = "доход"
	expenseType = "расход"
)

type Transaction struct {
	ID          int
	Type        string
	Amount      float64
	Category    string
	Description string
}

type FinanceTracker struct {
	transactions []Transaction
	nextID       int
}

func validateAmount(amount float64) error {
	if math.IsNaN(amount) || math.IsInf(amount, 0) {
		return errors.New("сумма должна быть конечным числом")
	}
	if amount <= 0 {
		return errors.New("сумма должна быть больше нуля")
	}
	return nil
}

func validateCategory(category string) error {
	if strings.TrimSpace(category) == "" {
		return errors.New("категория не должна быть пустой")
	}
	return nil
}

func (tracker *FinanceTracker) AddTransaction(transactionType string, amount float64, category, description string) error {
	if transactionType != incomeType && transactionType != expenseType {
		return errors.New("неизвестный тип операции")
	}
	if err := validateAmount(amount); err != nil {
		return err
	}
	if err := validateCategory(category); err != nil {
		return err
	}

	tracker.nextID++
	tracker.transactions = append(tracker.transactions, Transaction{
		ID:          tracker.nextID,
		Type:        transactionType,
		Amount:      amount,
		Category:    strings.TrimSpace(category),
		Description: strings.TrimSpace(description),
	})
	return nil
}

func (tracker FinanceTracker) Balance() (float64, float64, float64) {
	var income, expense float64
	for _, transaction := range tracker.transactions {
		if transaction.Type == incomeType {
			income += transaction.Amount
		} else {
			expense += transaction.Amount
		}
	}
	return income, expense, income - expense
}

func (tracker FinanceTracker) ExpensesByCategory() map[string]float64 {
	expenses := make(map[string]float64)
	for _, transaction := range tracker.transactions {
		if transaction.Type == expenseType {
			expenses[transaction.Category] += transaction.Amount
		}
	}
	return expenses
}

func (tracker FinanceTracker) FindByCategory(category string) []Transaction {
	category = strings.TrimSpace(category)
	var found []Transaction
	for _, transaction := range tracker.transactions {
		if transaction.Category == category {
			found = append(found, transaction)
		}
	}
	return found
}

func (tracker *FinanceTracker) DeleteTransaction(id int) error {
	for index, transaction := range tracker.transactions {
		if transaction.ID == id {
			tracker.transactions = append(tracker.transactions[:index], tracker.transactions[index+1:]...)
			return nil
		}
	}
	return fmt.Errorf("операция с идентификатором %d не найдена", id)
}
