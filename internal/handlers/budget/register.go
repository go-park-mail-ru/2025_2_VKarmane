package budget

import (
	"github.com/gorilla/mux"
)

func Register(r *mux.Router, uc BudgetUseCase) {
	// Бюджеты не используются фронтендом, поэтому регистрация пустая
}
