package l2

type Config struct {
	Endpoint     string `yaml:"endpoint"`
	BlockThreads uint64 `yaml:"block_threads"`
}
