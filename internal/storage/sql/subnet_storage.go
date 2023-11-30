package sqlstorage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/chirikova/go-anti-brute-force/internal/storage"
)

type SubNetStorage struct {
	connection *Storage
	table      string
}

func NewSubNetStorage(connection *Storage, table string) storage.SubNetStoragable {
	return &SubNetStorage{
		connection: connection,
		table:      table,
	}
}

func (s *SubNetStorage) Add(subnet string) error {
	query := fmt.Sprintf("INSERT INTO %s (subnet) VALUES($1)", s.table)

	_, err := s.connection.db.ExecContext(
		s.connection.ctx,
		query,
		subnet,
	)

	return err
}

func (s *SubNetStorage) Remove(subnet string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE subnet = $1", s.table)

	_, err := s.connection.db.ExecContext(
		s.connection.ctx,
		query,
		subnet,
	)

	return err
}

func (s *SubNetStorage) HasIP(subnet string) (bool, error) {
	query := fmt.Sprintf("SELECT id FROM %s WHERE $1::inet <<= subnet", s.table)
	result := make([]string, 0)
	err := s.connection.db.SelectContext(
		s.connection.ctx,
		&result,
		query,
		subnet,
	)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return true, err
	}

	return len(result) != 0, err
}
