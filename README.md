# DOMAIN CRUD GENERATOR

## How to run this program
1. Rename file `config.json.template` to `config.json`
2. Custom config
3. Run command below
```bash
go run main.go
```
4. Copy your result domain in folder result

### Config Description
- app_name : your service name
- domain_name : your domain name (separate with " **-** " if you need)
- table_name : your table name to query (separate with " **_** " like your real table name)
- attributes :
    - column_name : column in your table
    - type :  
        type | in golang
        --- | ---
        `int`, `integer` | int
        `int32` | int32
        `int64`, `bigint`, `int8`, `unix` | int64
        `string`, `varchar`, `text` | string
        `float32` | float32
        `float`, `float64`, `double` | float64
        `bool`, `boolean` | bool
        `[]int`, `int[]` | []int
        `[]int32` | []int32
        `[]int64`, `bigint[]`, `int8[]` | []int64
        `[]string`, `varchar[]` | []string
        `[]float32` | []float32
        `[]float64`, `double[]` | []float64
        `time`, `time.Time`, `timestamp`, `date` | time.Time
- options : fill it free
- redis_ttl & memcache_ttl : fill it if you using the cache

