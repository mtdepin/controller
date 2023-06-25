package config

type LoggerConfig struct {
	Level string
}

type JaegerConfig struct {
	Url string
}

type CommonNodeConfig struct {
	Id         string
	Name       string
	Node_group string
	Region     string
}

type MqConfig struct {
	Server             string
	Topic              string
	Topic_heartBeat    string
	Broadcast_group_id string
}

type RequestConfig struct {
	Max      int
	TimeOut  int
	Protocol string
}

type PprofConfig struct {
	pprof bool
}

type RedisConfig struct {
	LockAddress  string
	LockPassword string
	FidAddress   string
	FidPassword  string
}

type SearcherConfig struct {
	Url string
}

type StrategyConfig struct {
	Url string
}

type TaskTrackerConfig struct {
	Url string
}

type AccountConfig struct {
	Url string
}

type UserConfig struct {
	Url string
}

type SchedulerConfig struct {
	Url string
}

type MontiorConfig struct {
	AccessToken string
	Secret      string
}

type LockConfig struct {
	Prefix string
}

type PrometheusConfig struct {
	Url string
}

type SuperClusterConfig struct {
	Region string
}
