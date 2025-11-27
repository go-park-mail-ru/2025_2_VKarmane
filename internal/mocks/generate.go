// Package mocks provides gomock-generated mocks for testing
package mocks

//go:generate go run go.uber.org/mock/mockgen -destination=mock_finance_client.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto FinanceServiceClient
//go:generate go run go.uber.org/mock/mockgen -destination=mock_finance_service.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/usecase FinanceService
//go:generate go run go.uber.org/mock/mockgen -destination=mock_finance_usecase.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/grpc FinanceUseCase
//go:generate go run go.uber.org/mock/mockgen -destination=mock_finance_repository.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/service FinanceRepository

//go:generate go run go.uber.org/mock/mockgen -destination=mock_budget_client.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/proto BudgetServiceClient
//go:generate go run go.uber.org/mock/mockgen -destination=mock_budget_service.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/usecase BudgetService
//go:generate go run go.uber.org/mock/mockgen -destination=mock_budget_usecase.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/grpc BudgetUseCase
//go:generate go run go.uber.org/mock/mockgen -destination=mock_budget_repository.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/service BudgetRepository

//go:generate go run go.uber.org/mock/mockgen -destination=mock_auth_client.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/proto AuthServiceClient
//go:generate go run go.uber.org/mock/mockgen -destination=mock_auth_service.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/usecase AuthService
//go:generate go run go.uber.org/mock/mockgen -destination=mock_auth_usecase.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/grpc AuthUseCase
//go:generate go run go.uber.org/mock/mockgen -destination=mock_auth_repository.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/service AuthRepository

//go:generate go run go.uber.org/mock/mockgen -destination=mock_account_repository.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/account/service AccountRepository

//go:generate go run go.uber.org/mock/mockgen -destination=mock_operation_repository.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/operations/service OperationRepository
//go:generate go run go.uber.org/mock/mockgen -destination=mock_category_repository.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/category/service CategoryRepository

//go:generate go run go.uber.org/mock/mockgen -destination=mock_balance_service.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/account/service BalanceService

//go:generate go run go.uber.org/mock/mockgen -destination=mock_operation_service.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/operations/service OperationService
//go:generate go run go.uber.org/mock/mockgen -destination=mock_category_service.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/category/service CategoryService

//go:generate go run go.uber.org/mock/mockgen -destination=mock_image_service.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/image/service ImageService

//go:generate go run go.uber.org/mock/mockgen -destination=mock_balance_usecase.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/account/handlers BalanceUseCase

//go:generate go run go.uber.org/mock/mockgen -destination=mock_category_usecase.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/category/usecase CategoryUseCase

//go:generate go run go.uber.org/mock/mockgen -destination=mock_image_usecase.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/image/usecase ImageUseCase
//go:generate go run go.uber.org/mock/mockgen -destination=mock_operation_usecase.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/operations/handlers OperationUseCase

//go:generate go run go.uber.org/mock/mockgen -destination=mock_clock.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock Clock
//go:generate go run go.uber.org/mock/mockgen -destination=mock_kafka.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/kafka 
//go:generate go run go.uber.org/mock/mockgen -destination=mock_logger.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger Logger
//go:generate go run go.uber.org/mock/mockgen -destination=mock_image_storage.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/image/repository ImageStorage
