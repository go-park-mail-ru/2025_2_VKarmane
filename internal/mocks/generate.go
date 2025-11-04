// Package mocks provides gomock-generated mocks for testing
package mocks

//go:generate mockgen -destination=mock_user_repository.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/service UserRepository
//go:generate mockgen -destination=mock_account_repository.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/service AccountRepository
//go:generate mockgen -destination=mock_budget_repository.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/service BudgetRepository
//go:generate mockgen -destination=mock_operation_repository.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/service OperationRepository
//go:generate mockgen -destination=mock_category_repository.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/service CategoryRepository

//go:generate mockgen -destination=mock_auth_service.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/auth AuthService
//go:generate mockgen -destination=mock_balance_service.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/balance BalanceService
//go:generate mockgen -destination=mock_budget_service.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/budget BudgetService
//go:generate mockgen -destination=mock_operation_service.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/operation OperationService
//go:generate mockgen -destination=mock_category_service.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/category CategoryService
//go:generate mockgen -destination=mock_profile_service.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/profile ProfileService
//go:generate mockgen -destination=mock_image_service.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/service/image ImageService

//go:generate mockgen -destination=mock_auth_usecase.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/auth AuthUseCase
//go:generate mockgen -destination=mock_balance_usecase.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/balance BalanceUseCase
//go:generate mockgen -destination=mock_budget_usecase.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/handlers/budget BudgetUseCase
//go:generate mockgen -destination=mock_category_usecase.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/category CategoryUseCase
//go:generate mockgen -destination=mock_profile_usecase.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/profile ProfileUseCase
//go:generate mockgen -destination=mock_image_usecase.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase/image ImageUseCase

//go:generate mockgen -destination=mock_clock.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/clock Clock
//go:generate mockgen -destination=mock_logger.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger Logger
//go:generate mockgen -destination=mock_image_storage.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/repository/storage/image ImageStorage

