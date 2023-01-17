package factories

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/rs/zerolog"

	"github.com/koko1123/flow-go-1/module"
	builder "github.com/koko1123/flow-go-1/module/builder/collection"
	finalizer "github.com/koko1123/flow-go-1/module/finalizer/collection"
	"github.com/koko1123/flow-go-1/module/mempool"
	"github.com/koko1123/flow-go-1/network"
	"github.com/koko1123/flow-go-1/storage"
)

type BuilderFactory struct {
	db               *badger.DB
	mainChainHeaders storage.Headers
	trace            module.Tracer
	opts             []builder.Opt
	metrics          module.CollectionMetrics
	pusher           network.Engine // engine for pushing finalized collection to consensus committee
	log              zerolog.Logger
}

func NewBuilderFactory(
	db *badger.DB,
	mainChainHeaders storage.Headers,
	trace module.Tracer,
	metrics module.CollectionMetrics,
	pusher network.Engine,
	log zerolog.Logger,
	opts ...builder.Opt,
) (*BuilderFactory, error) {

	factory := &BuilderFactory{
		db:               db,
		mainChainHeaders: mainChainHeaders,
		trace:            trace,
		metrics:          metrics,
		pusher:           pusher,
		log:              log,
		opts:             opts,
	}
	return factory, nil
}

func (f *BuilderFactory) Create(
	clusterHeaders storage.Headers,
	clusterPayloads storage.ClusterPayloads,
	pool mempool.Transactions,
) (module.Builder, *finalizer.Finalizer, error) {

	build, err := builder.NewBuilder(
		f.db,
		f.trace,
		f.mainChainHeaders,
		clusterHeaders,
		clusterPayloads,
		pool,
		f.log,
		f.opts...,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create builder: %w", err)
	}

	final := finalizer.NewFinalizer(
		f.db,
		pool,
		f.pusher,
		f.metrics,
	)

	return build, final, nil
}
