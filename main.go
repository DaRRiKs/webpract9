package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const menu = `
===== PERSONAL FINANCE TRACKER =====
1. Добавить доход
2. Добавить расход
3. История операций
4. Текущий баланс
5. Расходы по категориям
6. Поиск по категории
7. Удалить операцию
0. Выход`

func main() {
	if err := run(); err != nil && !errors.Is(err, io.EOF) {
		fmt.Println("Ошибка чтения:", err)
	}
}

func run() error {
	reader := bufio.NewReader(os.Stdin)
	tracker := FinanceTracker{}
	for {
		fmt.Println(menu)
		choice, err := readLine(reader, "Выберите действие: ")
		if err != nil {
			return err
		}

		switch choice {
		case "1", "2":
			transactionType := incomeType
			if choice == "2" {
				transactionType = expenseType
			}
			err = addTransaction(reader, &tracker, transactionType)
		case "3":
			printTransactions(tracker.transactions, "История операций пуста.")
		case "4":
			income, expense, balance := tracker.Balance()
			fmt.Printf("Общие доходы: %.2f тг\n", income)
			fmt.Printf("Общие расходы: %.2f тг\n", expense)
			fmt.Printf("Текущий баланс: %.2f тг\n", balance)
		case "5":
			printExpenses(tracker.ExpensesByCategory())
		case "6":
			var category string
			category, err = readCategory(reader)
			if err == nil {
				printTransactions(tracker.FindByCategory(category), "Операций в выбранной категории нет.")
			}
		case "7":
			err = deleteTransaction(reader, &tracker)
		case "0":
			return nil
		default:
			err = errors.New("несуществующий пункт меню; выберите число от 0 до 7")
		}

		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			fmt.Println("Ошибка:", err)
		}
	}
}

func readLine(reader *bufio.Reader, prompt string) (string, error) {
	fmt.Print(prompt)
	line, err := reader.ReadString('\n')
	if errors.Is(err, io.EOF) && len(line) > 0 {
		err = nil
	}
	return strings.TrimSpace(line), err
}

func readCategory(reader *bufio.Reader) (string, error) {
	category, err := readLine(reader, "Категория: ")
	if err != nil {
		return "", err
	}
	return category, validateCategory(category)
}

func addTransaction(reader *bufio.Reader, tracker *FinanceTracker, transactionType string) error {
	input, err := readLine(reader, "Сумма: ")
	if err != nil {
		return err
	}
	amount, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return errors.New("некорректная сумма; введите число, например 4500 или 4500.50")
	}
	if err := validateAmount(amount); err != nil {
		return err
	}
	category, err := readCategory(reader)
	if err != nil {
		return err
	}
	description, err := readLine(reader, "Описание: ")
	if err != nil {
		return err
	}
	if err := tracker.AddTransaction(transactionType, amount, category, description); err != nil {
		return err
	}
	fmt.Println("Операция успешно добавлена!")
	return nil
}

func deleteTransaction(reader *bufio.Reader, tracker *FinanceTracker) error {
	input, err := readLine(reader, "Идентификатор: ")
	if err != nil {
		return err
	}
	id, err := strconv.Atoi(input)
	if err != nil {
		return errors.New("некорректный идентификатор; введите целое число")
	}
	if err := tracker.DeleteTransaction(id); err != nil {
		return err
	}
	fmt.Println("Операция успешно удалена!")
	return nil
}

func printTransactions(transactions []Transaction, emptyMessage string) {
	if len(transactions) == 0 {
		fmt.Println(emptyMessage)
		return
	}
	for _, transaction := range transactions {
		fmt.Printf("ID: %d | Тип: %s | Категория: %s | Сумма: %.2f тг\n",
			transaction.ID, transaction.Type, transaction.Category, transaction.Amount)
	}
}

func printExpenses(expenses map[string]float64) {
	if len(expenses) == 0 {
		fmt.Println("Расходов пока нет.")
		return
	}
	for category, amount := range expenses {
		fmt.Printf("%s: %.2f тг\n", category, amount)
	}
}
