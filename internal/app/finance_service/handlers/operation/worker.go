package operation

import (
	"context"
	"sync"

	finpb "github.com/go-park-mail-ru/2025_2_VKarmane/internal/app/finance_service/proto"
	"github.com/go-park-mail-ru/2025_2_VKarmane/internal/models"
)

func ProcessCSV(
	ctx context.Context,
	ops []models.CreateOperationRequest,
	userID, accountID int,
	client finpb.FinanceServiceClient,
) ([]models.OperationResponse, []error) {

	const workers = 10

	jobs := make(chan models.CreateOperationRequest)
	results := make(chan models.OperationResponse)
	errs := make(chan error)

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for req := range jobs {
				op, err := client.CreateOperation(ctx,
					CreateOperationRequestToProto(req, userID, accountID),
				)
				if err != nil {
					errs <- err
					continue
				}

				results <- ProtoOperationToResponse(op)
			}
		}()
	}

	go func() {
		for _, req := range ops {
			jobs <- req
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
		close(errs)
	}()

	var created []models.OperationResponse
	var errorsList []error

	finishedResults := false
    finishedErrors := false

    for !(finishedResults && finishedErrors) {
        select {
        case r, ok := <-results:
            if !ok {
                finishedResults = true
                continue
            }
            created = append(created, r)

        case e, ok := <-errs:
            if !ok {
                finishedErrors = true
                continue
            }
            errorsList = append(errorsList, e)
        }
    }

	return created, errorsList
}
