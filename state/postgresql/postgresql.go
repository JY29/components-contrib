/*
Copyright 2021 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package postgresql

import (
	"context"
	"reflect"

	"github.com/JY29/components-contrib/metadata"
	"github.com/JY29/components-contrib/state"
	"github.com/dapr/kit/logger"
)

// PostgreSQL state store.
type PostgreSQL struct {
	logger   logger.Logger
	dbaccess dbAccess
}

// NewPostgreSQLStateStore creates a new instance of PostgreSQL state store.
func NewPostgreSQLStateStore(logger logger.Logger) state.Store {
	dba := newPostgresDBAccess(logger)

	return newPostgreSQLStateStore(logger, dba)
}

// newPostgreSQLStateStore creates a newPostgreSQLStateStore instance of a PostgreSQL state store.
// This unexported constructor allows injecting a dbAccess instance for unit testing.
func newPostgreSQLStateStore(logger logger.Logger, dba dbAccess) *PostgreSQL {
	return &PostgreSQL{
		logger:   logger,
		dbaccess: dba,
	}
}

// Init initializes the SQL server state store.
func (p *PostgreSQL) Init(metadata state.Metadata) error {
	return p.dbaccess.Init(metadata)
}

// Features returns the features available in this state store.
func (p *PostgreSQL) Features() []state.Feature {
	return []state.Feature{state.FeatureETag, state.FeatureTransactional, state.FeatureQueryAPI}
}

// Delete removes an entity from the store.
func (p *PostgreSQL) Delete(ctx context.Context, req *state.DeleteRequest) error {
	return p.dbaccess.Delete(ctx, req)
}

// BulkDelete removes multiple entries from the store.
func (p *PostgreSQL) BulkDelete(ctx context.Context, req []state.DeleteRequest) error {
	return p.dbaccess.BulkDelete(ctx, req)
}

// Get returns an entity from store.
func (p *PostgreSQL) Get(ctx context.Context, req *state.GetRequest) (*state.GetResponse, error) {
	return p.dbaccess.Get(ctx, req)
}

// BulkGet performs a bulks get operations.
func (p *PostgreSQL) BulkGet(ctx context.Context, req []state.GetRequest) (bool, []state.BulkGetResponse, error) {
	// TODO: replace with ExecuteMulti for performance
	return false, nil, nil
}

// Set adds/updates an entity on store.
func (p *PostgreSQL) Set(ctx context.Context, req *state.SetRequest) error {
	return p.dbaccess.Set(ctx, req)
}

// BulkSet adds/updates multiple entities on store.
func (p *PostgreSQL) BulkSet(ctx context.Context, req []state.SetRequest) error {
	return p.dbaccess.BulkSet(ctx, req)
}

// Multi handles multiple transactions. Implements TransactionalStore.
func (p *PostgreSQL) Multi(ctx context.Context, request *state.TransactionalStateRequest) error {
	return p.dbaccess.ExecuteMulti(ctx, request)
}

// Query executes a query against store.
func (p *PostgreSQL) Query(ctx context.Context, req *state.QueryRequest) (*state.QueryResponse, error) {
	return p.dbaccess.Query(ctx, req)
}

// Close implements io.Closer.
func (p *PostgreSQL) Close() error {
	if p.dbaccess != nil {
		return p.dbaccess.Close()
	}
	return nil
}

// Returns the dbaccess property.
// This method is used in tests.
func (p *PostgreSQL) GetDBAccess() dbAccess {
	return p.dbaccess
}

func (p *PostgreSQL) GetComponentMetadata() map[string]string {
	metadataStruct := postgresMetadataStruct{}
	metadataInfo := map[string]string{}
	metadata.GetMetadataInfoFromStructType(reflect.TypeOf(metadataStruct), &metadataInfo)
	return metadataInfo
}
