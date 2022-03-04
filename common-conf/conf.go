package common_conf

type S3Config struct {
	Endpoint             string `json:"endp"`
	Region               string `json:"region,omitempty"`
	AKID                 string `json:"accesskey"`
	AKSK                 string `json:"secretkey"`
	Bucket               string `json:"bucket"`
	FileEncryptPassword  string `json:"encryption,omitempty"`
	UseTLS               bool   `json:"useTLS"`            // globally, both backend and final url will use this schema
	ReverseProxyEndPoint string `json:"urlendp,omitempty"` // if not set, not proxied, else change domain and port, start from http(s)
}

type WebConfig struct {
	ListenOn    string `json:"listen"`
	Environment string `json:"environ"` // prod OR dev
	DomainName  string `json:"domain"`
	BizName     string `json:"bizname"`
}

type ReCaptchaConfig struct {
	Secret     string  `json:"sitesecret"`
	SiteKey    string  `json:"sitekey"`
	Threshold  float32 `json:"threshold"`
	InChina    bool    `json:"china_user"`
	DomainName string  `json:"domain"`
}

type RedisConfig struct {
	DBNum    int    `json:"db_num,omitempty"`
	Addr     string `json:"addr"`
	Username string `json:"username,omitempty"`
	Password string `json:"password"`
}

type TLSConfig struct {
	CertPath string `json:"cert_path,omitempty"`
	KeyPath  string `json:"key_path,omitempty"`
}

type GRPCConfig struct {
	RegistryAddr      string `json:"registry_addr"`
	PublicKey         string `json:"pubkey"`            // ed25519, in base64
	PrivateKey        string `json:"privkey,omitempty"` // ed25519, in base64
	Role              int    `json:"role"`              // server = 0 or worker = 1
	HeartBeatInterval int    `json:"heartbeat"`
}

type WorkPoolConfig struct {
	MaxJobsInParallel int `json:"max_parallel"`
	MaxQueueSize      int `json:"max_queue"`
	MaxSingleNodeJobs int `json:"max_node"`
}
