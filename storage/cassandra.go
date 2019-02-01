package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"gitlab.corp.cloudsimple.com/cloudsimple/csos/incubator/shark-tank/types"
)

const (
	tableName = "example_table"
	keyspace  = "example"
)

type CassandraRepository struct {
	sess *gocql.Session
}

func NewCassandraRepository(address string, port int) (*CassandraRepository, error) {
	sess, err := createSession(address, port)
	if err != nil {
		return nil, err
	}
	if err := createKeyspaceAndTable(sess); err != nil {
		return nil, err
	}
	return &CassandraRepository{
		sess: sess,
	}, nil
}

func createSession(address string, port int) (*gocql.Session, error) {
	c := gocql.NewCluster(address)
	const protoVersion = 4
	const timeout = 15 * time.Second
	c.Port, c.Timeout, c.ConnectTimeout, c.ProtoVersion = port, timeout, timeout, protoVersion
	return c.CreateSession()
}

func (r *CassandraRepository) Create(ctx context.Context, entity *types.ExampleEntity) error {
	entity.ID = gocql.TimeUUID()
	q := r.sess.Query(fmt.Sprintf(`INSERT INTO %s.%s (id, message) VALUES (%s, '%s')`,
		keyspace, tableName, entity.ID, entity.Message)).WithContext(ctx)
	defer q.Release()
	return q.Exec()
}

func (r *CassandraRepository) Inject(ctx context.Context, entity *types.ExampleEntity) error {
	q := r.sess.Query(entity.Message).WithContext(ctx)
	defer q.Release()
	return q.Exec()
}

func (r *CassandraRepository) Get(ctx context.Context, id gocql.UUID) (*types.ExampleEntity, error) {
	q := r.sess.Query(fmt.Sprintf(`SELECT id, message FROM %s.%s WHERE id = %s`,
		keyspace, tableName, id.String())).WithContext(ctx)
	defer q.Release()

	entity := types.ExampleEntity{}
	if err := q.Scan(&entity.ID, &entity.Message); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *CassandraRepository) List(ctx context.Context) ([]types.ExampleEntity, error) {
	q := r.sess.Query(fmt.Sprintf(`SELECT id, message FROM %s.%s`,
		keyspace, tableName)).WithContext(ctx)
	defer q.Release()

	iter := q.Iter()
	entity := types.ExampleEntity{}
	list := make([]types.ExampleEntity, 0, iter.NumRows())
	for iter.Scan(&entity.ID, &entity.Message) {
		list = append(list, entity)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return list, nil
}

func createKeyspaceAndTable(sess *gocql.Session) error {
	keyspaceQuery := fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};`, keyspace)
	tableQuery := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.%s (id uuid PRIMARY KEY, message text );`, keyspace, tableName)

	for _, q := range []string{keyspaceQuery, tableQuery} {
		if err := sess.Query(q).Exec(); err != nil {
			return err
		}
	}
	return nil
}
