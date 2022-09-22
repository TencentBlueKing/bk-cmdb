## Business Topology Cache
  This package used to cache all the business's topology from the root nodes business all the way
to the lowest node module.
  It's used by various scenes, which is need to get business topology frequently and care about
the performance. like job's scheduled tasks.

  This business's topology cache has features as follows:
 - the  cache is a brief topology of this business, which contains the basic information with
  the object, instance id and name.
 - this cache is refreshed when a business's topology changed, such as a custom level instance is
  added, removed. or a set, module is added or removed. this is an event-drive mechanism, so that
  cache can be refreshed in time.
 - this cache has a ttl for several hours, which help us to clean the cache automatically when a
  business is deleted or archived.
 - all the cache refreshed every 15 minutes no matter event occurred or not. it's a safety
  mechanism to ensure the cache is correct.
 - if we cannot find business topology from the cache, we read it from the db directly.