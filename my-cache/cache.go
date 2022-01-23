package cache

type cache interface {
	Get()
	Set()
	Delete()
	New()
}
