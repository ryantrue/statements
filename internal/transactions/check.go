package transactions

import (
	"context"
	"statements/internal/database"
)

// existsTransaction проверяет, существует ли транзакция в базе данных (включая больше параметров)
func existsTransaction(accountNumber, date, debit, credit, documentNumber, paymentDescription, debitAccount, creditAccount, inn, name, innC, nameC string) (bool, error) {
	var exists bool
	err := database.DB.QueryRowContext(context.Background(),
		`SELECT EXISTS(
			SELECT 1 FROM transactions
			WHERE account_number = $1
			AND date = $2
			AND debit = $3
			AND credit = $4
			AND document_number = $5
			AND payment_description = $6
			AND debit_account = $7
			AND credit_account = $8
			AND inn = $9
			AND name = $10
			AND inn_c = $11
			AND name_c = $12
		)`,
		accountNumber, date, debit, credit, documentNumber, paymentDescription, debitAccount, creditAccount, inn, name, innC, nameC).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}
