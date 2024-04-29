package table

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
)

type Worker struct {
	Address common.Address `gorm:"column:address;type:bytea;not null;primaryKey;"`
	Network string         `gorm:"column:network;type:text;not null;primaryKey;"`
	Name    string         `gorm:"column:name;type:text;not null;primaryKey;"`
}

func (*Worker) TableName() string {
	return "node_worker"
}

func (w *Worker) Import(worker *schema.Worker) {
	w.Address = worker.Address
	w.Network = worker.Network
	w.Name = worker.Name
}

func (w *Worker) Export() *schema.Worker {
	return &schema.Worker{
		Address: w.Address,
		Network: w.Network,
		Name:    w.Name,
	}
}

type Workers []Worker

func (w *Workers) Export() []*schema.Worker {
	workers := make([]*schema.Worker, 0)

	for _, worker := range *w {
		exportedWorker := worker.Export()
		workers = append(workers, exportedWorker)
	}

	return workers
}

func (w *Workers) Import(workers []*schema.Worker) {
	*w = make([]Worker, 0, len(workers))

	for _, worker := range workers {
		var tWorker Worker

		tWorker.Import(worker)
		*w = append(*w, tWorker)
	}
}
