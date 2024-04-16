package generator

// type in go
const (
	intStr            string = "int"
	int32Str          string = "int32"
	int64Str          string = "int64"
	stringStr         string = "string"
	boolStr           string = "bool"
	float32Str        string = "float32"
	float64Str        string = "float64"
	arrayOfIntStr     string = "[]int"
	arrayOfInt32Str   string = "[]int32"
	arrayOfInt64Str   string = "[]int64"
	arrayOfStringStr  string = "[]string"
	arrayOfFloat32Str string = "[]float32"
	arrayOfFloat64Str string = "[]float64"
	timeStr           string = "time.Time"
)

// type data array
var mapIsDataArray = map[string]bool{
	arrayOfIntStr:     true,
	arrayOfInt32Str:   true,
	arrayOfInt64Str:   true,
	arrayOfStringStr:  true,
	arrayOfFloat32Str: true,
	arrayOfFloat64Str: true,
}

// nil data in DB
var mapNilDataDB = map[string]string{
	intStr:           `0`,
	int32Str:         `0`,
	int64Str:         `0`,
	stringStr:        `""`,
	boolStr:          `false`,
	arrayOfIntStr:    `ARRAY[]::int[]`,
	arrayOfInt32Str:  `ARRAY[]::int[]`,
	arrayOfInt64Str:  `ARRAY[]::bigint[]`,
	arrayOfStringStr: `ARRAY[]::varchar[]`,
	timeStr:          `'0001-01-01T00:00:00Z'::timestamp`,
}

// replacer
const (
	appNameReplacer                  string = "[appname]"
	domainNameLowerCaseReplacer      string = "[domainname]"
	domainNameCamelCaseReplacer      string = "[DomainName]"
	domainNameLowerCamelCaseReplacer string = "[domainName]"
	attributeReplacer                string = "[attribute]"
	attributeQueryReplacer           string = "[attributeQuery]"
	useCacheParamReplacer            string = "[useCacheParam]"
	acronymReplacer                  string = "[acronym]"
	cacheTTLReplacer                 string = "[cacheTTL]"
	memcacheTTLReplacer              string = "[memcacheTTL]"
	tableNameReplacer                string = "[tableName]"
	attributesListReplacer           string = "[attributesList]"
	paramNumberListReplacer          string = "[paramNumberList]"
	attributeSetListReplacer         string = "[attributeSetList]"
	idSetReplacer                    string = "[idSet]"
	scanGetAttributeReplacer         string = "[scanGetAttribute]"
	scanInsertAttributeReplacer      string = "[scanInsertAttribute]"
	scanUpdateAttributeReplacer      string = "[scanUpdateAttribute]"
	copyAttributeReplacer            string = "[copyAttribute]"
)

const (
	packageCode string = `package [domainname]`

	resourceTDKLib string = `"github.com/tokopedia/tdk/go/app/resource"`
	redisLib       string = `redislib "github.com/tokopedia/[appname]/lib/redislib"`
	memcacheLib    string = `"github.com/tokopedia/[appname]/pkg/lib/memcache"`
	cacheLib       string = `"github.com/tokopedia/[appname]/common/cache"`
)

var crudComponent []string = []string{"READ", "CREATE", "UPDATE", "DELETE"}

// format code
const (
	// [init] - import
	importInitCode string = `import (
    // FIX ME
    "context"
	"fmt"

)`

	// [init] - data type
	domainInterface   string = `[DomainName]DomainItf interface`
	domainStruct      string = `[DomainName]Domain struct`
	resourceInterface string = `[DomainName]ResourceItf interface`
	resourceStruct    string = `[DomainName]Resource struct`

	// [init] - resource struct
	dbResourceStruct       string = `DB resource.SQLDB`
	redisResourceStruct    string = `Redis redislib.RedisItf`
	memcacheResourceStruct string = `Memcache memcache.MemcacheItf`

	// [init] - domain struct
	resourceDomainStruct string = `resource [DomainName]ResourceItf`

	// parameter
	cacheParam string = `setCacheConfig cache.SetCacheConfig`

	// [init] - domain interface data
	domainInterfaceData string = `
    Get[DomainName]ByIDsBulk(ctx context.Context, [domainName]IDs []int64 [useCacheParam]) (map[int64][DomainName], error)
    Insert[DomainName](ctx context.Context, param Insert[DomainName]Param) (int64, error)
    Update[DomainName](ctx context.Context, param Update[DomainName]Param) error
    Delete[DomainName](ctx context.Context, [domainName]ID int64) error
    `

	// [init] - singleflight
	singleflightDomainInterfaceData string = `Get[DomainName]ByIDsBulkWithSingleFlight(ctx context.Context, [domainName]IDs []int64 [useCacheParam]) (map[int64][DomainName], error)`

	// [init] - resource interface data
	resourceInterfaceDataDB string = `
    // database
    get[DomainName]ByIDsBulkDB(ctx context.Context, [domainName]IDs []int64) (map[int64][DomainName], error)
    insert[DomainName]DB(ctx context.Context, param Insert[DomainName]Param) (int64, error)
    update[DomainName]DB(ctx context.Context, param Update[DomainName]Param) error
    delete[DomainName]DB(ctx context.Context, [domainName]ID int64) error
    `
	resourceInterfaceDataCache string = `
    // cache
    get[DomainName]ByIDsBulkCache(ctx context.Context, [domainName]IDs []int64) (map[int64][DomainName], error)
    set[DomainName]ByIDsBulkCache(ctx context.Context, map[DomainName] map[int64][DomainName], cacheDuration int) error
    delete[DomainName]ByIDCache(ctx context.Context, [domainName]ID int64) error
    `

	// [init] - func
	funcInitDomainCode string = `
func InitDomain(rsc [DomainName]ResourceItf) [DomainName]Domain {
	return [DomainName]Domain{
		resource: rsc,
	}
}`

	funcGetDomainWithCacheCode string = `
func ([domainName]Domain [DomainName]Domain) Get[DomainName]ByIDsBulk(ctx context.Context, [domainName]IDs []int64, setCacheConfig cache.SetCacheConfig) (map[int64][DomainName], error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "[domainname].Get[DomainName]ByIDsBulk")
	defer span.Finish()

	result := make(map[int64][DomainName])

	[domainName]IDs = util.ArrayInt64RemoveElement([domainName]IDs, 0)
	[domainName]IDs = util.UniqueArrayInt64([domainName]IDs)
	if len([domainName]IDs) == 0 {
		return nil, ers.ErrorAddTrace(ers.ErrorInvalidParameter)
	}

	var (
		db[DomainName]s   map[int64][DomainName]
		notCachedIDs []int64
		notFoundIDs  []int64
		cacheResults map[int64][DomainName] = make(map[int64][DomainName])
		err          error
	)

	if !setCacheConfig.IsSkipGetCache {
		cacheResults, err = [domainName]Domain.resource.get[DomainName]ByIDsBulkCache(ctx, [domainName]IDs)
		if err != nil {
			return nil, ers.ErrorAddTrace(err)
		}
	}

	cached[DomainName]s := make(map[int64][DomainName])
	for _, [domainName]ID := range [domainName]IDs {
		if cacheResult, ok := cacheResults[[domainName]ID]; ok && cacheResult.ID != 0 {
			cached[DomainName]s[[domainName]ID] = cacheResult
		} else {
			notCachedIDs = append(notCachedIDs, [domainName]ID)
		}
	}

	defer func() {
		if len(db[DomainName]s) > 0 {
			cacheDuration := CacheTTL
			if setCacheConfig.CacheDuration > 0 {
				cacheDuration = setCacheConfig.CacheDuration
			}
			if err = [domainName]Domain.resource.set[DomainName]ByIDsBulkCache(ctx, db[DomainName]s, cacheDuration); err != nil {
				contextlib.PrintErrorCtx(ctx, err)
			}
		}
	}()

	if len(notCachedIDs) > 0 {
		db[DomainName]s, err = [domainName]Domain.resource.get[DomainName]ByIDsBulkDB(ctx, notCachedIDs)
		if err != nil {
			return nil, ers.ErrorAddTrace(err)
		}
	}

	for _, [domainName]ID := range [domainName]IDs {
		[domainName] := [DomainName]{}
		if val, ok := cached[DomainName]s[[domainName]ID]; ok {
			[domainName] = val
		} else if val, ok := db[DomainName]s[[domainName]ID]; ok {
			[domainName] = val
		} else {
			notFoundIDs = append(notFoundIDs, [domainName]ID)
		}

		if [domainName].ID > 0 {
			result[[domainName]ID] = [domainName]
		}
	}

	if len(notFoundIDs) > 0 && featureflag.IsActive(featureflag.PrintDebugErrorDomainFunction) {
		contextlib.PrintErrorCtx(
			ctx,
			ers.ErrorAddTrace(ers.[DomainName]NotFound),
			fmt.Sprintf(
				"[Get[DomainName]ByIDsBulk] [domainName]ID: %+v|cacheResults: %+v|db[DomainName]s: %+v",
				notFoundIDs,
				cacheResults,
				db[DomainName]s,
			),
		)
	}

	return result, nil
}`

	funcGetDomainWithoutCacheCode string = `
func ([domainName]Domain [DomainName]Domain) Get[DomainName]ByIDsBulk(ctx context.Context, [domainName]IDs []int64) (map[int64][DomainName], error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "[domainname].Get[DomainName]ByIDsBulk")
	defer span.Finish()

	result := make(map[int64][DomainName])

	[domainName]IDs = util.ArrayInt64RemoveElement([domainName]IDs, 0)
	[domainName]IDs = util.UniqueArrayInt64([domainName]IDs)
	if len([domainName]IDs) == 0 {
		return nil, ers.ErrorAddTrace(ers.ErrorInvalidParameter)
	}

	var (
		db[DomainName]s  map[int64][DomainName]
		notFoundIDs []int64
		err         error
	)

	db[DomainName]s, err = [domainName]Domain.resource.get[DomainName]ByIDsBulkDB(ctx, [domainName]IDs)
	if err != nil {
		return nil, ers.ErrorAddTrace(err)
	}

	for _, [domainName]ID := range [domainName]IDs {
		[domainName] := [DomainName]{}
		if val, ok := db[DomainName]s[[domainName]ID]; ok {
			[domainName] = val
		} else {
			notFoundIDs = append(notFoundIDs, [domainName]ID)
		}
		
		result[[domainName]ID] = [domainName]
	}

	if len(notFoundIDs) > 0 && featureflag.IsActive(featureflag.PrintDebugErrorDomainFunction) {
		contextlib.PrintErrorCtx(
			ctx,
			ers.ErrorAddTrace(ers.[DomainName]NotFound),
			fmt.Sprintf(
				"[Get[DomainName]ByIDsBulk] [domainName]ID: %+v|db[DomainName]s: %+v",
				notFoundIDs,
				db[DomainName]s,
			),
		)
	}

	return result, nil
}`

	funcInsertDomainCode string = `
func ([domainName]Domain [DomainName]Domain) Insert[DomainName](ctx context.Context, param Insert[DomainName]Param) (id int64, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "[domainname].Insert[DomainName]")
	defer span.Finish()

	id, err = [domainName]Domain.resource.insert[DomainName]DB(ctx, param)
	if err != nil {
		return 0, ers.ErrorAddTrace(err)
	}
	return id, nil
}`

	funcUpdateDomainCode string = `
func ([domainName]Domain [DomainName]Domain) Update[DomainName](ctx context.Context, param Update[DomainName]Param) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "[domainname].Update[DomainName]")
	defer span.Finish()

	err = [domainName]Domain.resource.update[DomainName]DB(ctx, param)
	if err != nil {
		return ers.ErrorAddTrace(err)
	}

	err = [domainName]Domain.resource.delete[DomainName]ByIDCache(ctx, param.ID)
	if err != nil {
		return ers.ErrorAddTrace(err)
	}

	return nil
}`

	funcDeleteDomainCode string = `
func ([domainName]Domain [DomainName]Domain) Delete[DomainName](ctx context.Context, id int64) (err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "[domainname].Delete[DomainName]")
	defer span.Finish()

	err = [domainName]Domain.resource.delete[DomainName]DB(ctx, id)
	if err != nil {
		return ers.ErrorAddTrace(err)
	}

	err = [domainName]Domain.resource.delete[DomainName]ByIDCache(ctx, id)
	if err != nil {
		return ers.ErrorAddTrace(err)
	}

	return nil
}`

	// [type] - struct format
	structRead string = `
type [DomainName] struct {
    [attribute]
}`

	structCreate string = `
type Insert[DomainName]Param struct {
    [attribute]
}`

	structUpdate string = `
type Update[DomainName]Param struct {
    [attribute]
}`

	structDelete string = `
type Delete[DomainName]Param struct {
    [attribute]
}`

	structSingleflight string = `
type SF[DomainName]Response struct {
    [DomainName]s map[int64][DomainName]
	SFTraceID     string
}`

	// [const] - cache
	constCache string = `
const (
    CacheKey[DomainName]ByID string = "[acronym]:[domainname]:id:%d"
    CacheTTL int = [cacheTTL]
    MemcacheTTL int = [memcacheTTL]
)
    `
	// [const] - singleflight
	constSingleflight string = `
const (
    sfKeyGet[DomainName]ByIDsBulk string = "get[domainname]byidsbulk"
    sfTraceKey string = "domain[DomainName]TraceID"
)
    `

	// [const] - query
	constReadQuery string = "QueryGet[DomainName] string = `" + `
    SELECT
        [attributeQuery]
	FROM [tableName]
    ` + "`\n" +
		"QueryGet[DomainName]ByIDsBulk string = QueryGet[DomainName] + ` WHERE [tableName].id IN ($arr)`"

	constCreateQuery string = "QueryInsert[DomainName] string = `" + `
    INSERT INTO [tableName]
    ([attributesList])
    VALUES([paramNumberList])
	RETURNING id
    ` + "`"

	constUpdateQuery string = "QueryUpdate[DomainName] string = `" + `
    UPDATE [tableName]
    SET
        [attributeSetList]
	WHERE [idSet]
    ` + "`\n"

	constDeleteQuery string = "QueryDelete[DomainName] string = `" + `
    DELETE FROM [tableName]
	WHERE [idSet]
    ` + "`\n"

	// [database] - import
	importDatabaseCode string = `import (
	// FIX ME
	"context"
	"time"

	"github.com/lib/pq"
	contextlib "github.com/tokopedia/[appname]/common/context"
	dbutil "github.com/tokopedia/[appname]/common/database"
	ers "github.com/tokopedia/[appname]/common/error"
	opentracing "github.com/tokopedia/tdk/go/tracer/v2"
	)`

	// [database] - all func
	funcDatabaseCode string = `
func (rsc [DomainName]Resource) get[DomainName]ByIDBulkDB(ctx context.Context, [domainName]IDs []int64) (map[int64][DomainName], error) {
	span, ctx := opentracing.StartExternalSpanFromContext(ctx, "[domainname].db.get[DomainName]ByIDsBulkDB", opentracing.WithSQLSpan(QueryGet[DomainName]ByIDsBulk, map[string]interface{}{}))
	defer span.Finish()
	
	queryTimeoutDuration, errFF := dbutil.GetQueryTimeoutDuration()
	if errFF != nil {
		contextlib.PrintErrorCtx(ctx, errFF, "[Feature Flag] QueryTimeoutDuration")
	}
	ctx, cancel := context.WithTimeout(ctx, time.Duration(queryTimeoutDuration)*time.Millisecond)
	defer cancel()

	result := make(map[int64][DomainName])

	query := dbutil.ReplaceInParam(1, len([domainName]IDs), QueryGet[DomainName]ByIDsBulk)
	var queryParam []interface{}
	for _, [domainName]ID := range [domainName]IDs {
		queryParam = append(queryParam, [domainName]ID)
	}

	rows, err := rsc.DB.QueryContext(ctx, query, queryParam...)
	if err != nil {
		return nil, ers.ErrorAddTrace(err)
	}
	defer rows.Close()

	for rows.Next() {
		var data [DomainName]
		err := rows.Scan(
			[scanGetAttribute]
		)

		if err != nil {
			return nil, ers.ErrorAddTrace(err)
		}

		result[data.ID] = data
	}

	if err = rows.Err(); err != nil {
		contextlib.PrintErrorCtx(ctx, errFF, "[ErrorRow]", err)
		return nil, ers.ErrorAddTrace(err)
	}
	
	return result, nil
}
	
func (rsc [DomainName]Resource) insert[DomainName]DB(ctx context.Context, param Insert[DomainName]Param) (id int64, err error) {
	span, ctx := opentracing.StartExternalSpanFromContext(ctx, "[domainname].db.insert[DomainName]", opentracing.WithSQLSpan(QueryInsert[DomainName], map[string]interface{}{}))
	defer span.Finish()

	db := rsc.DB.GetMaster()
	err = db.QueryRowContext(
		ctx,
		QueryInsert[DomainName],
		[scanInsertAttribute]
	).Scan(&id)
	if err != nil {
		return 0, ers.ErrorAddTrace(err)
	}

	return id, nil
}

func (rsc [DomainName]Resource) update[DomainName]DB(ctx context.Context, param Update[DomainName]Param) (err error) {
	span, ctx := opentracing.StartExternalSpanFromContext(ctx, "[domainname].db.update[DomainName]", opentracing.WithSQLSpan(QueryUpdate[DomainName], map[string]interface{}{}))
	defer span.Finish()

	db := rsc.DB.GetMaster()
	_, err = db.ExecContext(
		ctx,
		QueryUpdate[DomainName],
		[scanUpdateAttribute]
	)
	if err != nil {
		return ers.ErrorAddTrace(err)
	}

	return nil
}

func (rsc [DomainName]Resource) delete[DomainName]DB(ctx context.Context, id int64) (err error) {
	span, ctx := opentracing.StartExternalSpanFromContext(ctx, "[domainname].db.delete[DomainName]", opentracing.WithSQLSpan(QueryDelete[DomainName], map[string]interface{}{}))
	defer span.Finish()

	db := rsc.DB.GetMaster()
	_, err = db.ExecContext(
		ctx,
		QueryDelete[DomainName],
		id,
	)
	if err != nil {
		return ers.ErrorAddTrace(err)
	}

	return nil
}`

	// [cache] - import
	importCacheCode string = `import (
	// FIX ME
	"context"
	"fmt"

	contextlib "github.com/tokopedia/[appname]/common/context"
	ers "github.com/tokopedia/[appname]/common/error"
	opentracing "github.com/tokopedia/tdk/go/tracer/v2"
)`

	// [cache] - all func
	funcCacheCode string = `
func constructCacheKey[DomainName]ByID([domainName]ID int64) string {
	return fmt.Sprintf(CacheKey[DomainName]ByID, [domainName]ID)
}

func (rsc [DomainName]Resource) get[DomainName]ByIDsBulkCache(ctx context.Context, [domainName]IDs []int64) (map[int64][DomainName], error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "[domainname].cache.get[DomainName]ByIDsBulkCache")
	defer span.Finish()

	result := make(map[int64][DomainName])
	notMapcacheIDs := []int64{}

	vals := make(map[string]interface{})
	for _, [domainName]ID := range [domainName]IDs {
		key := constructCacheKey[DomainName]ByID([domainName]ID)

		mapCacheData, ok := rsc.get[DomainName]ByIDMapCache(ctx, key)
		if ok {
			result[[domainName]ID] = mapCacheData
			continue
		}

		// set param for get redis
		vals[key] = &[DomainName]{}
		notMapcacheIDs = append(notMapcacheIDs, [domainName]ID)
	}

	if len(notMapcacheIDs) <= 0 {
		return result, nil
	}

	if err := rsc.Redis.GetBulkStructs(ctx, vals); err != nil {
		return map[int64][DomainName]{}, ers.ErrorAddTrace(err)
	}

	for _, notMapcacheID := range notMapcacheIDs {
		key := constructCacheKey[DomainName]ByID(notMapcacheID)
		itf := vals[key]
		data, _ := itf.(*[DomainName])
		result[notMapcacheID] = *data

		if data != nil && data.ID > 0 {
			rsc.Memcache.Set(ctx, key, *data, MemcacheTTL)
		}
	}

	return result, nil
}

func (rsc [DomainName]Resource) get[DomainName]ByIDMapCache(ctx context.Context, key string) ([DomainName], bool) {
	var data [DomainName]

	err := rsc.Memcache.Get(ctx, key, &data)
	if err != nil {
		contextlib.DebugPrintlnCtx(ctx, err, "[get[DomainName]ByIDMapCache] error")
		return [DomainName]{}, false
	}

	return data, true
}

func (rsc [DomainName]Resource) set[DomainName]ByIDsBulkCache(ctx context.Context, map[DomainName] map[int64][DomainName], cacheDuration int) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "[domainname].cache.set[DomainName]ByIDsBulkCache")
	defer span.Finish()

	vals := make(map[string]interface{})
	for [domainName]ID, [domainName] := range map[DomainName] {
		key := constructCacheKey[DomainName]ByID([domainName]ID)

		vals[key] = [domainName]
		rsc.Memcache.Set(ctx, key, [domainName], MemcacheTTL)
	}

	if err := rsc.Redis.SetBulkStructs(ctx, vals, cacheDuration); err != nil {
		return ers.ErrorAddTrace(err)
	}

	return nil
}

func (rsc [DomainName]Resource) delete[DomainName]Cache(ctx context.Context, id int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "[domiannname].cache.delete[DomainName]Cache")
	defer span.Finish()

	_, err := rsc.Redis.Delete(ctx, constructCacheKey[DomainName]ByID(id))
	if err != nil {
		return ers.ErrorAddTrace(err)
	}

	return nil
}`

	// [singleflight] - import
	importSingleflightCode string = `import (
	// FIX ME
	"context"
	"fmt"
	"sort"

	"github.com/tokopedia/[appname]/common/cache"
	contextlib "github.com/tokopedia/[appname]/common/context"
	ers "github.com/tokopedia/[appname]/common/error"
	opentracing "github.com/tokopedia/tdk/go/tracer/v2"
	"golang.org/x/sync/singleflight"
)`

	// [singleflight] - all func
	funcSingleflightCode string = `
var sf[DomainName] singleflight.Group

func ([acronym]Domain [DomainName]Domain) Get[DomainName]ByIDsBulkWithSingleFlight(ctx context.Context, [domainName]IDs []int64, status int, setCacheConfig cache.SetCacheConfig) ([][DomainName], error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "[domainname].Get[DomainName]ByIDsBulkWithSingleFlight")
	defer span.Finish()

	key := generateSFKeyBy[DomainName]IDs(sfKeyGet[DomainName]ByIDsBulk, [domainName]IDs)

	itf, err, _ := sf[DomainName].Do(key, func() (interface{}, error) {
		ctx, sfTraceID := contextlib.GenerateTraceID(ctx, sfTraceKey)
		[DomainName]s, err := [acronym]Domain.Get[DomainName]ByIDsBulk(ctx, [domainName]IDs, status, setCacheConfig)

		return SF[DomainName]Response{
			[DomainName]s: [DomainName]s,
			SFTraceID: sfTraceID,
		}, err
	})
	if err != nil {
		return nil, ers.ErrorAddTrace(err)
	}

	result, ok := itf.(SF[DomainName]Response)
	if !ok {
		return nil, ers.ErrorAddTrace(ers.ParseDataError)
	}

	opentracing.AddAttribute(ctx, sfTraceKey, result.SFTraceID)
	return copy[DomainName]s(result.[DomainName]s), nil
}

func generateSFKeyBy[DomainName]IDs(methodName string, [domainName]IDs []int64) string {
	sort.Slice([domainName]IDs, func(i, j int) bool {
		return [domainName]IDs[i] < [domainName]IDs[j]
	})

	return fmt.Sprintf(
		"%s:%v",
		methodName,
		[domainName]IDs,
	)
}

func copy[DomainName]s([domainName]s [][DomainName]) [][DomainName] {
	copyResult := make([][DomainName], len([domainName]s))
	for i, [domainName] := range [domainName]s {
		[copyAttribute]
	}

	return copyResult
}
	
`
)
