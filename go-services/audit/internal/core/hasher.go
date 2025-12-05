package core

type Hasher interface {
	Hash(prevHash string, log *LogEntry) string
}
