package data

type Opt func(*DB) error

func WithSrc(src Src) Opt {
	return func(db *DB) error {
		db.Src = src
		return nil
	}
}

func WithName(name string) Opt {
	return func(db *DB) error {
		db.Name = name
		return nil
	}
}

func WithHare(path string) Opt {
	return func(db *DB) error {
		h, err := OpenHare(path)
		if err != nil {
			return err
		}
		db.onDisk = true
		db.Database = h
		return nil
	}
}
