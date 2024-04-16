package table

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/rss3-network/global-indexer/schema"
)

type Indexer struct {
	Address common.Address `gorm:"column:address;primaryKey"`
	Network string         `gorm:"column:network;primaryKey"`
	Worker  string         `gorm:"column:worker;primaryKey"`
}

func (*Indexer) TableName() string {
	return "node_indexer"
}

func (i *Indexer) Import(indexer *schema.Indexer) (err error) {
	i.Address = indexer.Address
	i.Network = indexer.Network
	i.Worker = indexer.Worker

	return nil
}

func (i *Indexer) Export() (*schema.Indexer, error) {
	indexer := schema.Indexer{
		Address: i.Address,
		Network: i.Network,
		Worker:  i.Worker,
	}

	return &indexer, nil
}

type Indexers []Indexer

func (i *Indexers) Export() ([]*schema.Indexer, error) {
	indexers := make([]*schema.Indexer, 0)

	for _, indexer := range *i {
		exportedIndexer, err := indexer.Export()
		if err != nil {
			return nil, err
		}

		indexers = append(indexers, exportedIndexer)
	}

	return indexers, nil
}

func (i *Indexers) Import(indexers []*schema.Indexer) (err error) {
	*i = make([]Indexer, 0, len(indexers))

	for _, indexer := range indexers {
		var tIndexer Indexer

		if err = tIndexer.Import(indexer); err != nil {
			return err
		}

		*i = append(*i, tIndexer)
	}

	return nil
}
