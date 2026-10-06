# Caching

Caching is a fundamental part of working with remote data across the Internet,
where access is severely impacted by latency and transfer time. In real-world scenarios,
you'll always need to keep all the data you need close and sync the changes with the remote
CMS when they happen.

Gocontentful supports caching out of the box. The client can maintain a cache of an entire space
or a subset of the content types that can be initialized with a single method call:

```go
contentTypes := []string{"person", "pet"}
err = cc.UpdateCache(context, contentTypes, true)
```

This makes sense for client modes `ClientModeCDA` and `ClientModeCPA` and not for the management API.
The client will download all the entries, convert and store them in the case as
native Go value objects. This makes subsequent accesses to the space data an in-memory operation removing all the HTTP
overhead.

The first parameter is the context. If you don't use a context in your application or service just pass _context.Background()_

The third parameter of UpdateCache toggles asset caching on or off. If you deal with assets you want this to be always on.

## Full cache init and rebuild

By default the client will cache the whole space using 4 parallel workers to speed up the process.
This is safe since Contentful allows up to 5 concurrent connections.
If you have content types that have a lot of entries, it might make sense to keep them close to each other
in the content types slice passed to UpdateCache(), so that they will run in parallel and not one after the other.

All gocontentful functions that query the space cache-transparent: if a cache is available data will be loaded from
there, otherwise it will be sourced from Contentful. This doesn't apply to _GetFilteredXYZ()_ calls that
always need to pass the query to Contentful.

Gocontentful also supports selective entry and asset cache updates through the following method:

```go
err = cc.UpdateCacheForEntity(context, sysType, contentType, entityID string)
```

When something changes in the space at Contentful you need to update the cache. For this to happen you need to set
up a webhook at Contentful and handle its calls in your service through a public HTTP listener.
When a webhook call gets in, you have the choice of updating your cache in different ways:

- You can regenerate the entire CDA cache when something is published because you want production data to
  be 100% up to date in your application. This can get slow and expensive.
- You can alternatively update a single entry in the cache. This is usually the case for the CPA cache because
  it's a lot faster and that works well for preview features.
- You can use the Sync API, but only limited to `ClientModeCDA`, as explained in the following paragraph.

In any case, if an update fails UpdateCache returns the error and the previous cache is preserved to prevent service disruption.
In the unfortunate case a service or application needs to start and Contentful is not available, Gocontentful can work
in an offline mode if you call _SetOfflineFallback_ on the client after you create it passing the path to a space export file.

The gocontentful API can work entirely offline too. In this case a cache is created from a space export file and most of the
features are available (pretty obviously, those that don't require live access to the space, like custom queries). If you update
the export file periodically you can even update the cache from the updated file.

## Entries with older versions and cache resets

Before a full rebuild replaces the cache, every entry Contentful returns is compared with the cached copy.
If Contentful returns an older published version (`sys.publishedVersion`) than the one in cache, the cached copy is kept,
together with its references, and the rest of the cache is refreshed. Entries and assets the kept copy links to (through
reference fields, asset fields or rich text) are kept from the previous cache as well when the fresh data lacks them, so
that its references keep resolving; they are listed in `result.RetainedDependencies`. This protects the cache from stale responses
served by the Contentful API. The update still succeeds: call _UpdateCacheWithResult_ to know which entries were retained:

```go
result, err := cc.UpdateCacheWithResult(ctx, contentTypes, true)
if err != nil {
	// the update failed and the previous cache is unchanged
}
if result.Degraded() {
	// result.Retained lists the entries kept from the previous cache
}
```

Each retained entry is also logged as a warning with its cached and incoming versions, followed by a summary,
and _ClientStats_ reports how many entries are currently retained.

Entries stay retained until Contentful serves a version at least as new as the cached one. That never happens when the
older versions are legitimate, for instance after an environment was re-created from another one. In that case call
_ForceUpdateCache_: it rebuilds the cache and replaces it without comparing versions, listing the overwritten entries in
`result.Replaced`. The rebuild still has to succeed, and a snapshot that is empty or holds less than a tenth of the cached
entries is refused, leaving the cache as it was. _ForceUpdateCache_ is not available in sync mode.

Rebuilds are serialized with entity updates; an update waiting for another one gives up when its context ends or the rebuild times out. A request arriving while a rebuild is queued joins it and gets its result,
and a _ForceUpdateCache_ request turns the queued rebuild into a reset. The content types and asset caching requested by the
first _UpdateCache_ call apply to all later rebuilds.

## Sync API support

In versions v1.0.12 and newer, gocontentful supports the Contentful Sync API and that's now the recommended way to cache spaces and manage updates.
Sync is enabled by default when you create a client with CDA mode.
To enable or disable support for the Sync API explicitly, you can call the SetSyncMode method on the client:

```go
cc.SetSyncMode(true)
```

With sync on, the cache updates will happen transparently through downloads of incremental changes.
The syntax to update the cache doesn't change, just call _UpdateCache_ on the client as usual.

The initialization of the cache will be slower when _SyncMode_ is on compared to the legacy full cache init because sync calls cannot be parallelized.
Subsequent updates though will be much faster because only changes in the space from the previous sync will be downloaded.
This includes entries and assets that were deleted. In case of need you can call _ResetSync()_ to start over from a fresh empty cache.

Note that the Sync API is not officially supported by Contentful on the Preview API. At the time of this writing it seems to work but use it at your own risk.

## Cache timeout

Cache update operations time out by default after 120 seconds. This makes sure that no
routine is left hanging, blocking subsequent updates in case the main application or service
recovers from a panic. If you need to increase this limit because you have a huge space with
a lot of entries you can use the _SetCacheUpdateTimeout_ method. See the [API Reference](./04-api-reference) for details.
A rebuild that times out returns an error, and a caller never waits longer than the timeout for a rebuild, even when
the Contentful API asks to retry later: it gets a timeout error while the rebuild goes on in the background. A caller
whose context is canceled stops waiting and gets the context error, while the rebuild goes on for any other caller
waiting on it.

## Asset caching

If you use assets in your space, then you absolutely need to enable them in the _UpdateCache_ call.
Otherwise, every time an entry needs to resolve a reference to an asset that single asset will be downloaded
and that for large spaces with thousands of assets can lead to incredibly slow operation.

## When to use and not use caching

Simple answer is: you should almost always use caching. The only scenario where not using
a cache on the client is better is when you only need to download a very limited amount
of entries (in the order of less than some hundreds) and do that at significant distance in time
(e.g. every hour). In this case your application code can be simpler and there won't be any
performance penalty. The other case is when you need to run a lot of custom queries or
use XPath, which is currently not supported by gocontentful directly.
