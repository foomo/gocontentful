package testapi

import (
	"errors"
)

type CacheStats struct {
	AssetCount   int
	ContentTypes []string
	EntryCount   int
	ParentCount  int
}

func (cc *ContentfulClient) GetCacheStats() (*CacheStats, error) {
	if cc == nil {
		return nil, errors.New("GetCacheStats: no client available")
	}
	return &CacheStats{
		AssetCount:   len(cc.Cache.assets),
		ContentTypes: cc.Cache.contentTypes,
		EntryCount:   len(cc.Cache.idContentTypeMap),
		ParentCount:  len(cc.Cache.parentMap),
	}, nil
}

func (cc *ContentfulClient) SetProductInCache(product *CfProduct) {
	if cc.Cache == nil {
		return
	}
	cc.cacheMutex.productGcLock.Lock()
	defer cc.cacheMutex.productGcLock.Unlock()
	cc.cacheMutex.idContentTypeMapGcLock.Lock()
	defer cc.cacheMutex.idContentTypeMapGcLock.Unlock()
	cc.cacheMutex.parentMapGcLock.Lock()
	defer cc.cacheMutex.parentMapGcLock.Unlock()
	cc.Cache.entryMaps.product[product.Sys.ID] = product
	cc.Cache.idContentTypeMap[product.Sys.ID] = product.Sys.ContentType.Sys.ID
}
