package testapi

import "errors"

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
	cc.Cache.entryMaps.productGcLock.Lock()
	defer cc.Cache.entryMaps.productGcLock.Unlock()
	cc.Cache.idContentTypeMapGcLock.Lock()
	defer cc.Cache.idContentTypeMapGcLock.Unlock()
	cc.Cache.parentMapGcLock.Lock()
	defer cc.Cache.parentMapGcLock.Unlock()
	cc.Cache.entryMaps.product[product.Sys.ID] = product
	cc.Cache.idContentTypeMap[product.Sys.ID] = product.Sys.ContentType.Sys.ID
	return
}
