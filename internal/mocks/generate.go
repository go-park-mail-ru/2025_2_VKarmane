// Package mocks provides gomock-generated mocks for testing
package mocks

//go:generate go run go.uber.org/mock/mockgen -destination=mock_finance_client.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto FinanceServiceClient
//go:generate go run go.uber.org/mock/mockgen -destination=mock_finance_service.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/usecase FinanceService
//go:generate go run go.uber.org/mock/mockgen -destination=mock_finance_usecase.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/grpc FinanceUseCase
//go:generate go run go.uber.org/mock/mockgen -destination=mock_finance_repository.go -package=mocks github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/service FinanceRepository

