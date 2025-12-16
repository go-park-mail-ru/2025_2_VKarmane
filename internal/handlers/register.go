package handlers

import (
	"github.com/gorilla/mux"

	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/handlers/auth"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/handlers/profile"
	authpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/auth_service/proto"
	budget "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/handlers"
	bdgpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/budget_service/proto"
	balance "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/handlers/account"
	category "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/handlers/category"
	operation "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/handlers/operation"
	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	image "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/image/handlers"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/logger"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/usecase"
	kafkautils "github.com/go-park-mail-ru/2025_2_VKarmane/internal/utils/kafka"
)

type Registrator struct {
	uc  *usecase.UseCase
	log logger.Logger
}

func NewRegistrator(uc *usecase.UseCase, log logger.Logger) *Registrator {
	return &Registrator{
		uc:  uc,
		log: log,
	}

}

<<<<<<< HEAD
func (r *Registrator) RegisterAll(publicRouter *mux.Router, protectedRouter *mux.Router, uc *usecase.UseCase, log logger.Logger, authClient authpb.AuthServiceClient, budgetClient bdgpb.BudgetServiceClient, finClient finpb.FinanceServiceClient, kafkaProducer kafkautils.KafkaProducer) {
	auth.Register(publicRouter, protectedRouter, log, authClient)
	balance.Register(protectedRouter, finClient)
	budget.Register(protectedRouter, budgetClient)
	operation.Register(protectedRouter, finClient, uc.ImageUC, kafkaProducer)
	category.Register(protectedRouter, finClient, uc.ImageUC, kafkaProducer)
=======
func (r *Registrator) RegisterAll(publicRouter *mux.Router, protectedRouter *mux.Router, uc *usecase.UseCase, log logger.Logger, authClient authpb.AuthServiceClient, budgetClient bdgpb.BudgetServiceClient, finClient finpb.FinanceServiceClient) {
	auth.Register(publicRouter, protectedRouter, log, authClient)
	balance.Register(protectedRouter, finClient)
	budget.Register(protectedRouter, budgetClient)
	operation.Register(protectedRouter, finClient, uc.ImageUC)
	category.Register(protectedRouter, finClient, uc.ImageUC)
>>>>>>> main
	profile.Register(protectedRouter, uc.ImageUC, authClient)
	image.Register(protectedRouter, uc.ImageUC)
}
