package db

type Opt func(*DB)

func WithUID(uid string) Opt {
	return func(db *DB) {
		db.UID = uid
	}
}

func WithName(name string) Opt {
	return func(db *DB) {
		db.Name = name
	}
}
