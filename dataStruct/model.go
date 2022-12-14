package dataStruct

type GlobalConfig struct {
	ChainName        string           `yaml:"chain_name"`
	Version          string           `yaml:"version"`
	Client           ClientConfig     `yaml:"client"`
	Common           CommonConfig     `yaml:"common"`
	Cache            CacheConfig      `yaml:"cache"`
	Consensus        ConsensusConfig  `yaml:"consensus"`
	Block            BlockConfig      `yaml:"block"`
	P2p              P2pConfig        `yaml:"p2p"`
	Consul           ConsulConfig     `yaml:"ConsulConfig"`
	EtcdPrefixConfig EtcdPrefixConfig `yaml:"EtcdKeyPrefix"`
}
type EtcdPrefixConfig struct {
	UserOperationPrefix string `yaml:"user_operation"`
	UserInfoPrefix      string `yaml:"user_info"`
}
type ClientConfig struct {
	Organization string       `yaml:"organization"`
	CryptoConfig CryptoConfig `yaml:"crypto_config"`
}

type CryptoConfig struct {
	CryptoStore string `yaml:"crypto_store"`
	Path        string `yaml:"path"`
}

type CommonConfig struct {
	LogConfig    LogConfig               `yaml:"log_config"`
	CacheCluster []string                `yaml:"cache_cluster,flow"`
	LedgerName   map[string]LedgerConfig `yaml:"ledger_name"`
}

type LogConfig struct {
	RootPath     string `yaml:"root_path"`
	InfoLogName  string `yaml:"info_log_name"`
	DebugLogName string `yaml:"debug_log_name"`
	ErrorLogName string `yaml:"error_log_name"`
	LogLevel     string `yaml:"log_level"`
	OutputFile   bool   `yaml:"output_file"`
}

type LedgerConfig struct {
	Leader   string   `yaml:"leader"`
	Follower []string `yaml:"follower"`
}

type CacheConfig struct {
	CommonConfig RedisCommonConfig      `yaml:"common_config"`
	RedisGroup   map[string]RedisConfig `yaml:"redis_group"`
}

type RedisCommonConfig struct {
	DB            int    `yaml:"db"`
	Password      string `yaml:"password"`
	SyncInternal  int    `yaml:"sync_interval"`
	SyncSizeLimit int    `yaml:"sync_size_limit"`
	ExpireTime    int    `yaml:"expire_time"`
	Connection    int    `yaml:"connection"`
	Response      int    `yaml:"response"`
}

type RedisConfig struct {
	Host       string           `yaml:"host"`
	Port       string           `yaml:"port"`
	WebService WebServiceConfig `yaml:"webservice"`
}

type WebServiceConfig struct {
	URL string `yaml:"url"`
}

type ConsensusConfig struct {
	CommonConfig ConsensusCommonConfig `yaml:"common_config"`
	EtcdGroup    map[string]EtcdConfig `yaml:"etcd_group"`
	AsyncTask    map[string]TaskConfig `yaml:"async_task"`
}

type ConsensusCommonConfig struct {
	Timeout     int `yaml:"timeout"`
	MaxLiveDays int `yaml:"max_live_days"`
}

type EtcdConfig struct {
	HraftGrpcAddress string `yaml:"hraft_grpc_address"`
	BlockGrpcAddress string `yaml:"block_grpc_address"`
	Host             string `yaml:"host"`
	Port             string `yaml:"port"`
}

type TaskConfig struct {
	GenisisBlock GenisisBlockConfig `yaml:"genisis_block"`
	MinBlock     MinBlockConfig     `yaml:"min_block"`
	RehanceBlock RehanceBlockConfig `yaml:"rehance_block"`
	DailyBlock   DailyBlockConfig   `yaml:"daily_block"`
	BlockHeader  BlockHeaderConfig  `yaml:"block_header"`
}

type GenisisBlockConfig struct {
	Interval int `yaml:"interval"`
}

type MinBlockConfig struct {
	PackInterval int `yaml:"pack_interval"`
}

type RehanceBlockConfig struct {
	PackInterval int `yaml:"pack_interval"`
}

type DailyBlockConfig struct {
	PackInterval int `yaml:"pack_interval"`
}
type BlockHeaderConfig struct {
	Interval int `yaml:"interval"`
}

type BlockConfig struct {
	TDengineConfig  TDengineConfig  `yaml:"tdengine_config"`
	BlockFileConfig BlockFileConfig `yaml:"block_file_config"`
}

type TDengineConfig struct {
	MaxLiveDays int    `yaml:"max_live_days"`
	Hostname    string `yaml:"host_name"`
	Port        int    `yaml:"port"`
	User        string `yaml:"user"`
	Passwd      string `yaml:"passwd"`
	Keep        int    `yaml:"keep"`
	Driver      string `yaml:"driver"`
	DBName      string `yaml:"db_name"`
}

type BlockFileConfig struct {
	RootPath string `yaml:"root_path"`
}

type P2pConfig struct {
	Local LocalConfig `yaml:"local_config"`
}

type LocalConfig struct {
	Rendezvous string   `yaml:"rendezvous"`
	Pid        string   `yaml:"pid"`
	Port       string   `yaml:"port"`
	Host       string   `yaml:"host"`
	Groups     []string `yaml:"groups,flow"`
}

type ConsulConfig struct {
	ConsulAddress    string `yaml:"consul_address"`
	ConsulPort       string `yaml:"consul_port"`
	ID               string `yaml:"ID"`
	Name             string `yaml:"Name"`
	LocalAddress     string `yaml:"LocalAddress"`
	LocalServicePort string `yaml:"LocalServicePort"`
	HealthCheckID    string `yaml:"HealthCheckID"`
	HealthTCP        string `yaml:"HealthTCP"`
	HealthTimeout    string `yaml:"HealthTimeout"`
	HealthInterval   string `yaml:"HealthInterval"`
}

type OnlineConsulServer struct {
	ID            string `yaml:"ID"`
	ServerAddress string `yaml:"ServerAddress"`
	ServicePort   int    `yaml:"ServicePort"`
}

// ????????????
type UserInfo struct {
	UserName string `json:"UserName"`
	Password string `json:"Password"`
	Phone    string `json:"Phone"`
	Email    string `json:"Email"`
}
