package conf

type KafkaConf struct {
	Addr []string `yaml:"addr"`
}

type EtcdConf struct {
	Endpoint []string `yaml:"endpoint"`
	Key      string   `yaml:"key"`
}

type EsConf struct {
	Address []string `yaml:"address"`
}

type Conf struct {
	Kafka KafkaConf `yaml:"kafka"`
	Etcd  EtcdConf  `yaml:"etcd"`
	Es    EsConf    `yaml:"es"`
}
