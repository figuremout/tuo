package iorm

// TODO according to gorm
// put accumulator in
type DB struct {
	*Config
	Error        error
	RowsAffected int64
	Statement    *Statement
	// contains filtered or unexported fields
}

func Open() (*DB, error) {
	return &DB{}, nil
}

func (db *DB) Create() *DB {
	return nil
}

func (db *DB) Range() *DB {
	return nil
}

func (db *DB) Filter() *DB {
	return nil
}

func (db *DB) Keep() *DB {
	return nil
}

func (db *DB) Last() *DB {
	return nil
}
