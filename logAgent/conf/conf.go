package conf

type KafkaConf struct {
	Addr        []string `yaml:"addr"`
	ChanMaxSize int      `ymal:"chanMaxSize"`
}

type EtcdConf struct {
	Endpoint []string `yaml:"endpoint"`
	Key      string   `yaml:"key"`
}

type TailConf struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

type Conf struct {
	Kafka KafkaConf `yaml:"kafka"`
	Etcd  EtcdConf  `yaml:"etcd"`
}
