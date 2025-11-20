// Package mocks provides gomock-generated mocks for testing
package mocks

//go:generate go run go.uber.org/mock/mockgen -destination=mock_budget_client.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/proto BudgetServiceClient
//go:generate go run go.uber.org/mock/mockgen -destination=mock_budget_service.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/usecase BudgetService
//go:generate go run go.uber.org/mock/mockgen -destination=mock_budget_usecase.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/grpc BudgetUseCase
//go:generate go run go.uber.org/mock/mockgen -destination=mock_budget_repository.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/budget_service/service BudgetRepository

//go:generate go run go.uber.org/mock/mockgen -destination=mock_auth_client.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/proto AuthServiceClient
//go:generate go run go.uber.org/mock/mockgen -destination=mock_auth_service.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/usecase AuthService
//go:generate go run go.uber.org/mock/mockgen -destination=mock_auth_usecase.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/grpc AuthUseCase
//go:generate go run go.uber.org/mock/mockgen -destination=mock_auth_repository.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/auth_service/service AuthRepository

//go:generate go run go.uber.org/mock/mockgen -destination=mock_account_repository.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/service AccountRepository

//go:generate go run go.uber.org/mock/mockgen -destination=mock_operation_repository.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/service OperationRepository
//go:generate go run go.uber.org/mock/mockgen -destination=mock_category_repository.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/service CategoryRepository

//go:generate go run go.uber.org/mock/mockgen -destination=mock_balance_service.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/balance BalanceService

//go:generate go run go.uber.org/mock/mockgen -destination=mock_operation_service.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/operation OperationService
//go:generate go run go.uber.org/mock/mockgen -destination=mock_category_service.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/category CategoryService

//go:generate go run go.uber.org/mock/mockgen -destination=mock_image_service.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/image ImageService

//go:generate go run go.uber.org/mock/mockgen -destination=mock_balance_usecase.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/balance BalanceUseCase

//go:generate go run go.uber.org/mock/mockgen -destination=mock_category_usecase.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/category CategoryUseCase

//go:generate go run go.uber.org/mock/mockgen -destination=mock_image_usecase.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/image ImageUseCase
//go:generate go run go.uber.org/mock/mockgen -destination=mock_operation_usecase.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/operation OperationUseCase

//go:generate go run go.uber.org/mock/mockgen -destination=mock_clock.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock Clock
//go:generate go run go.uber.org/mock/mockgen -destination=mock_logger.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger Logger
//go:generate go run go.uber.org/mock/mockgen -destination=mock_image_storage.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/storage/image ImageStorage
