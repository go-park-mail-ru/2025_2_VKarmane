package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/account"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/dto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/user"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils"
	_ "github.com/lib/pq"
)

func main() {
	// Получаем параметры подключения из переменных окружения или используем значения по умолчанию
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5433")
	dbUser := getEnv("DB_USER", "vkarmane")
	dbPassword := getEnv("DB_PASSWORD", "vkarmane_password")
	dbName := getEnv("DB_NAME", "vkarmane")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("✅ Connected to database")

	ctx := context.Background()

	userRepo := user.NewPostgresRepository(db)
	accountRepo := account.NewPostgresRepository(db)

	// Очищаем тестовых пользователей перед созданием новых
	log.Println("🧹 Cleaning up existing test users...")
	testLogins := []string{"testuser1", "testuser2", "testuser3", "testuser4", "testuser5"}

	for _, login := range testLogins {
		// Получаем пользователя по логину
		userDB, err := userRepo.GetUserByLogin(ctx, login)
		if err != nil {
			// Пользователь не найден - пропускаем
			continue
		}

		log.Printf("  🗑️  Deleting test user: %s (ID: %d)", login, userDB.ID)

		// Получаем все счета пользователя
		accounts, err := accountRepo.GetAccountsByUser(ctx, userDB.ID)
		if err == nil && len(accounts) > 0 {
			// Удаляем связи пользователя со счетами (sharings)
			for _, acc := range accounts {
				_, err = db.ExecContext(ctx, `
					DELETE FROM sharings 
					WHERE user_id = $1 AND account_id = $2
				`, userDB.ID, acc.ID)
				if err != nil {
					log.Printf("    ⚠️  Failed to delete sharing for account %d: %v", acc.ID, err)
				}
			}

			// Удаляем счета пользователя (если они больше никому не принадлежат)
			for _, acc := range accounts {
				// Проверяем, есть ли еще связи с этим счетом
				var count int
				err = db.QueryRowContext(ctx, `
					SELECT COUNT(*) FROM sharings WHERE account_id = $1
				`, acc.ID).Scan(&count)
				if err == nil && count == 0 {
					// Если больше нет связей, удаляем счет
					_, err = db.ExecContext(ctx, `DELETE FROM account WHERE _id = $1`, acc.ID)
					if err != nil {
						log.Printf("    ⚠️  Failed to delete account %d: %v", acc.ID, err)
					} else {
						log.Printf("    ✅ Deleted account %d", acc.ID)
					}
				}
			}
		}

		// Удаляем пользователя (CASCADE удалит связанные записи)
		_, err = db.ExecContext(ctx, `DELETE FROM "user" WHERE _id = $1`, userDB.ID)
		if err != nil {
			log.Printf("    ⚠️  Failed to delete user %s: %v", login, err)
		} else {
			log.Printf("    ✅ Deleted user %s", login)
		}
	}

	log.Println("✅ Cleanup completed")

	// Тестовые пользователи
	testUsers := []struct {
		login     string
		email     string
		password  string
		firstName string
		lastName  string
		accounts  []struct {
			balance     float64
			accountType string
		}
	}{
		{
			login:     "testuser1",
			email:     "testuser1@example.com",
			password:  "password123",
			firstName: "Иван",
			lastName:  "Иванов",
			accounts: []struct {
				balance     float64
				accountType string
			}{
				{balance: 10000.50, accountType: "default"},
				{balance: 5000.00, accountType: "savings"},
			},
		},
		{
			login:     "testuser2",
			email:     "testuser2@example.com",
			password:  "password123",
			firstName: "Мария",
			lastName:  "Петрова",
			accounts: []struct {
				balance     float64
				accountType string
			}{
				{balance: 25000.75, accountType: "default"},
				{balance: 15000.00, accountType: "savings"},
				{balance: 5000.25, accountType: "investment"},
			},
		},
		{
			login:     "testuser3",
			email:     "testuser3@example.com",
			password:  "password123",
			firstName: "Алексей",
			lastName:  "Сидоров",
			accounts: []struct {
				balance     float64
				accountType string
			}{
				{balance: 5000.00, accountType: "default"},
			},
		},
		{
			login:     "testuser4",
			email:     "testuser4@example.com",
			password:  "password123",
			firstName: "Елена",
			lastName:  "Козлова",
			accounts: []struct {
				balance     float64
				accountType string
			}{
				{balance: 30000.00, accountType: "default"},
				{balance: 20000.00, accountType: "savings"},
			},
		},
		{
			login:     "testuser5",
			email:     "testuser5@example.com",
			password:  "password123",
			firstName: "Дмитрий",
			lastName:  "Смирнов",
			accounts: []struct {
				balance     float64
				accountType string
			}{
				{balance: 15000.50, accountType: "default"},
				{balance: 10000.00, accountType: "savings"},
				{balance: 7500.00, accountType: "investment"},
			},
		},
	}

	// Проверяем наличие валют в БД, если нет - создаем
	var currencyID int
	err = db.QueryRowContext(ctx, "SELECT _id FROM currency LIMIT 1").Scan(&currencyID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("⚠️  No currencies found in database, creating default currency (RUB)...")
			// Создаем валюту RUB по умолчанию
			err = db.QueryRowContext(ctx, `
				INSERT INTO currency (code, currency_name, logo_hashed_id, created_at)
				VALUES ('RUB', 'Российский рубль', 'e3b0c44298fc1c149afbf4c8996fb924', NOW())
				RETURNING _id
			`).Scan(&currencyID)
			if err != nil {
				log.Printf("❌ Failed to create default currency: %v", err)
				log.Fatalf("Cannot proceed without currency. Please create currency manually.")
			}
			log.Printf("✅ Created default currency RUB (ID: %d)", currencyID)
		} else {
			log.Printf("❌ Failed to check currencies: %v", err)
			log.Fatalf("Cannot proceed without currency. Please check database connection.")
		}
	} else {
		log.Printf("✅ Using existing currency ID: %d", currencyID)
	}

	for _, userData := range testUsers {
		// Хешируем пароль
		hashedPassword, err := utils.HashPassword(userData.password)
		if err != nil {
			log.Printf("❌ Failed to hash password for %s: %v", userData.login, err)
			continue
		}

		// Создаем пользователя
		userDB := dto.UserDB{
			FirstName:    userData.firstName,
			LastName:     userData.lastName,
			Email:        userData.email,
			Login:        userData.login,
			Password:     hashedPassword,
			Description:  fmt.Sprintf("Тестовый пользователь %s", userData.login),
			LogoHashedID: "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b",
		}

		userID, err := userRepo.CreateUser(ctx, userDB)
		if err != nil {
			log.Printf("❌ Failed to create user %s: %v", userData.login, err)
			// Возможно, пользователь уже существует - попробуем получить его ID
			if userDB, err := userRepo.GetUserByLogin(ctx, userData.login); err == nil {
				userID = userDB.ID
				log.Printf("ℹ️  User %s already exists, using existing ID: %d", userData.login, userID)
			} else {
				continue
			}
		} else {
			log.Printf("✅ Created user: %s (ID: %d)", userData.login, userID)
		}

		// Создаем счета для пользователя
		for _, accountData := range userData.accounts {
			now := time.Now()
			accountDB := account.AccountDB{
				Balance:    accountData.balance,
				Type:       accountData.accountType,
				CurrencyID: currencyID,
				CreatedAt:  now,
				UpdatedAt:  now,
			}

			accountID, err := accountRepo.CreateAccount(ctx, accountDB)
			if err != nil {
				log.Printf("❌ Failed to create account for user %s: %v", userData.login, err)
				continue
			}

			// Связываем счет с пользователем через sharings
			userAccountDB := account.UserAccountDB{
				UserID:    userID,
				AccountID: accountID,
				CreatedAt: now,
				UpdatedAt: now,
			}

			err = accountRepo.CreateUserAccount(ctx, userAccountDB)
			if err != nil {
				log.Printf("❌ Failed to link account %d to user %s: %v", accountID, userData.login, err)
				continue
			}

			log.Printf("  ✅ Created account %d (balance: %.2f, type: %s) for user %s",
				accountID, accountData.balance, accountData.accountType, userData.login)
		}
	}

	log.Println("✅ Seed completed successfully!")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
