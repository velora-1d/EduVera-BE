package client_temporal_inbound_adapter

import (
	"go.temporal.io/sdk/worker"

	"eduvera/internal/domain"
	"eduvera/internal/model"
	inbound_port "eduvera/internal/port/inbound"
	"eduvera/utils/activity"
	"eduvera/utils/log"
	"eduvera/utils/temporal"
)

type clientAdapter struct {
	domain domain.Domain
}

func NewClientAdapter(
	domain domain.Domain,
) inbound_port.ClientWorkflowPort {
	return &clientAdapter{
		domain: domain,
	}
}

func (a *clientAdapter) Upsert() {
	ctx := activity.NewContext("upsert_client_worker")

	w, err := temporal.NewWorker(ctx, model.UpsertClientWorkflowName)
	if err != nil {
		log.WithContext(ctx).Error("Unable to create worker", err)
		return
	}

	workflow := NewClientWorkflow(a.domain)

	w.RegisterWorkflow(workflow.UpsertClientWorkflow)
	w.RegisterActivity(a.domain.Client())

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.WithContext(ctx).Error("Unable to start worker", err)
		return
	}
}
