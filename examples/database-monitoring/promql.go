// Package monitoring demonstrates database observability patterns.
//
// This example shows PostgreSQL and Redis monitoring with shared PromQL
// expressions for alerts and dashboards.
package monitoring

import "github.com/lex00/wetwire-observability-go/promql"

// PostgreSQL Expressions

// PGConnectionsUsedExpr calculates the percentage of PostgreSQL connections in use.
var PGConnectionsUsedExpr = promql.Div(
	promql.Metric("pg_stat_activity_count"),
	promql.Metric("pg_settings_max_connections"),
)

// PGSlowQueriesExpr counts queries running longer than 30 seconds.
var PGSlowQueriesExpr = promql.Sum(
	promql.Vector("pg_stat_activity_max_tx_duration",
		promql.Match("state", "active")),
).By("datname")

// PGReplicationLagExpr calculates replication lag in seconds.
var PGReplicationLagExpr = promql.Max(
	promql.Metric("pg_replication_lag"),
).By("instance")

// PGDeadlocksExpr calculates the rate of deadlocks per database.
var PGDeadlocksExpr = promql.Sum(
	promql.Rate(promql.RangeVector("pg_stat_database_deadlocks", "5m")),
).By("datname")

// PGCacheHitRatioExpr calculates the buffer cache hit ratio.
var PGCacheHitRatioExpr = promql.Div(
	promql.Sum(promql.Metric("pg_stat_database_blks_hit")).By("datname"),
	promql.Add(
		promql.Sum(promql.Metric("pg_stat_database_blks_hit")).By("datname"),
		promql.Sum(promql.Metric("pg_stat_database_blks_read")).By("datname"),
	),
)

// Redis Expressions

// RedisConnectedClientsExpr returns the number of connected clients.
var RedisConnectedClientsExpr = promql.Metric("redis_connected_clients")

// RedisMemoryUsedExpr calculates Redis memory usage percentage.
var RedisMemoryUsedExpr = promql.Div(
	promql.Metric("redis_memory_used_bytes"),
	promql.Metric("redis_memory_max_bytes"),
)

// RedisKeysExpiringExpr calculates the rate of keys expiring.
var RedisKeysExpiringExpr = promql.Sum(
	promql.Rate(promql.RangeVector("redis_expired_keys_total", "5m")),
).By("instance")

// RedisEvictionsExpr calculates the rate of key evictions.
var RedisEvictionsExpr = promql.Sum(
	promql.Rate(promql.RangeVector("redis_evicted_keys_total", "5m")),
).By("instance")

// RedisCommandsPerSecExpr calculates commands processed per second.
var RedisCommandsPerSecExpr = promql.Sum(
	promql.Rate(promql.RangeVector("redis_commands_processed_total", "5m")),
).By("instance")

// RedisHitRatioExpr calculates the keyspace hit ratio.
var RedisHitRatioExpr = promql.Div(
	promql.Rate(promql.RangeVector("redis_keyspace_hits_total", "5m")),
	promql.Add(
		promql.Rate(promql.RangeVector("redis_keyspace_hits_total", "5m")),
		promql.Rate(promql.RangeVector("redis_keyspace_misses_total", "5m")),
	),
)
