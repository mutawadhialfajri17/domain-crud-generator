package generator

type RequestData struct {
	AppName     string      `json:"app_name"`
	DomainName  string      `json:"domain_name"`
	TableName   string      `json:"table_name"`
	Attributes  []Attribute `json:"attributes"`
	Options     Options     `json:"options"`
	RedisTTL    int         `json:"redis_ttl"`
	MemcacheTTL int         `json:"memcache_ttl"`
}

type Attribute struct {
	ColumnName string `json:"column_name"`
	Type       string `json:"type"`
}

type Options struct {
	IsUseRedisAndMemcache bool `json:"is_use_redis_and_memcache"`
	IsUseSingleflight     bool `json:"is_use_singleflight"`
}
