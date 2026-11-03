package ports

type Hasher interface {
	Hash(password string) (string, error)
	Check(password, hash string) (bool, error)
}
