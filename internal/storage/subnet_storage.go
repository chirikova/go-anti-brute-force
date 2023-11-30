package storage

type SubNetStoragable interface {
	Add(subnet string) error
	Remove(subnet string) error
	HasIP(subnet string) (bool, error)
}
