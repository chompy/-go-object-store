package storage

import (
	"gitlab.com/contextualcode/storage-backend/pkg/types"
)

// SqliteBackend handles connection to Sqlite database.
type SqliteBackend struct {
	Path string
	//	db   *sql.DB
}

/*func (b *SqliteBackend) init() error {
	var err error
	b.db, err = sql.Open("sqlite3", fmt.Sprintf("file:%s", b.Path))

	b.db.Query(
		`
		CREATE TABLE IF NOT EXISTS objects
		(
			uid			CHAR(32)	PRIMARY KEY		NOT NULL,
			created 	DATETIME	NOT NULL,
			modified	DATETIME	NOT NULL,
			type		CHAR(32)	NOT NULL
		)
		`,
	)

	b.db.Query(
		`
		CREATE TABLE IF NOT EXISTS values
		(
			id			INT			PRIMARY KEY		NOT NULL,
			ouid		CHAR(32)	NOT NULL,
			value		TEXT
		)
		`,
	)

	return errors.WithStack(err)
}*/

// Put uploads given object.
func (b *SqliteBackend) Put(o *types.Object) error {

	return nil
}

// Get downloads object.
func (b *SqliteBackend) Get(uid string) (*types.Object, error) {
	return nil, nil
}
