package monitoring

import (
	"github.com/lex00/wetwire-observability-go/grafana"
	"github.com/lex00/wetwire-observability-go/promql"
)

// PostgreSQL Panels

// PGConnectionUsagePanel displays connection pool usage.
var PGConnectionUsagePanel = grafana.Stat("Connection Usage").
	WithTargets(grafana.PromTargetExpr(PGConnectionsUsedExpr).WithRefID("A")).
	WithUnit(grafana.UnitPercentUnit)

// PGCacheHitPanel displays cache hit ratio.
var PGCacheHitPanel = grafana.Stat("Cache Hit Ratio").
	WithTargets(grafana.PromTargetExpr(PGCacheHitRatioExpr).WithRefID("A")).
	WithUnit(grafana.UnitPercentUnit)

// PGReplicationLagPanel displays replication lag over time.
var PGReplicationLagPanel = grafana.TimeSeries("Replication Lag").
	WithTargets(grafana.PromTargetExpr(
		promql.Max(promql.Metric("pg_replication_lag")).By("instance"),
	).WithRefID("A").WithLegendFormat("{{ instance }}")).
	WithUnit(grafana.UnitSeconds)

// PGDeadlocksPanel displays deadlock rate over time.
var PGDeadlocksPanel = grafana.TimeSeries("Deadlocks").
	WithTargets(grafana.PromTargetExpr(
		promql.Sum(promql.Rate(promql.RangeVector("pg_stat_database_deadlocks", "$__rate_interval"))).By("datname"),
	).WithRefID("A").WithLegendFormat("{{ datname }}")).
	WithUnit(grafana.UnitShort)

// PGActiveConnectionsPanel displays active connections over time.
var PGActiveConnectionsPanel = grafana.TimeSeries("Active Connections").
	WithTargets(grafana.PromTargetExpr(
		promql.Sum(promql.Vector("pg_stat_activity_count", promql.Match("state", "active"))).By("datname"),
	).WithRefID("A").WithLegendFormat("{{ datname }}")).
	WithUnit(grafana.UnitShort)

// Redis Panels

// RedisMemoryUsagePanel displays memory usage.
var RedisMemoryUsagePanel = grafana.Stat("Memory Usage").
	WithTargets(grafana.PromTargetExpr(RedisMemoryUsedExpr).WithRefID("A")).
	WithUnit(grafana.UnitPercentUnit)

// RedisHitRatioPanel displays cache hit ratio.
var RedisHitRatioPanel = grafana.Stat("Hit Ratio").
	WithTargets(grafana.PromTargetExpr(RedisHitRatioExpr).WithRefID("A")).
	WithUnit(grafana.UnitPercentUnit)

// RedisConnectedClientsPanel displays connected clients over time.
var RedisConnectedClientsPanel = grafana.TimeSeries("Connected Clients").
	WithTargets(grafana.PromTargetExpr(
		promql.Sum(promql.Metric("redis_connected_clients")).By("instance"),
	).WithRefID("A").WithLegendFormat("{{ instance }}")).
	WithUnit(grafana.UnitShort)

// RedisCommandsPanel displays commands per second.
var RedisCommandsPanel = grafana.TimeSeries("Commands/sec").
	WithTargets(grafana.PromTargetExpr(
		promql.Sum(promql.Rate(promql.RangeVector("redis_commands_processed_total", "$__rate_interval"))).By("instance"),
	).WithRefID("A").WithLegendFormat("{{ instance }}")).
	WithUnit(grafana.UnitShort)

// RedisMemoryTimeSeriesPanel displays memory usage over time.
var RedisMemoryTimeSeriesPanel = grafana.TimeSeries("Memory Usage").
	WithTargets(
		grafana.PromTargetExpr(
			promql.Sum(promql.Metric("redis_memory_used_bytes")).By("instance"),
		).WithRefID("A").WithLegendFormat("used"),
		grafana.PromTargetExpr(
			promql.Sum(promql.Metric("redis_memory_max_bytes")).By("instance"),
		).WithRefID("B").WithLegendFormat("max"),
	).
	WithUnit(grafana.UnitBytes)

// RedisEvictionsPanel displays evictions over time.
var RedisEvictionsPanel = grafana.TimeSeries("Key Evictions").
	WithTargets(grafana.PromTargetExpr(
		promql.Sum(promql.Rate(promql.RangeVector("redis_evicted_keys_total", "$__rate_interval"))).By("instance"),
	).WithRefID("A").WithLegendFormat("{{ instance }}")).
	WithUnit(grafana.UnitShort)

// Dashboard Definitions

// PostgresDashboard is the main PostgreSQL monitoring dashboard.
var PostgresDashboard = grafana.NewDashboard("postgres-metrics", "PostgreSQL Metrics").
	WithTags("postgres", "database").
	WithRows(
		grafana.NewRow("Overview").WithPanels(
			PGConnectionUsagePanel,
			PGCacheHitPanel,
		),
		grafana.NewRow("Connections").WithPanels(
			PGActiveConnectionsPanel,
		),
		grafana.NewRow("Replication & Health").WithPanels(
			PGReplicationLagPanel,
			PGDeadlocksPanel,
		),
	)

// RedisDashboard is the main Redis monitoring dashboard.
var RedisDashboard = grafana.NewDashboard("redis-metrics", "Redis Metrics").
	WithTags("redis", "database", "cache").
	WithRows(
		grafana.NewRow("Overview").WithPanels(
			RedisMemoryUsagePanel,
			RedisHitRatioPanel,
		),
		grafana.NewRow("Clients & Commands").WithPanels(
			RedisConnectedClientsPanel,
			RedisCommandsPanel,
		),
		grafana.NewRow("Memory & Evictions").WithPanels(
			RedisMemoryTimeSeriesPanel,
			RedisEvictionsPanel,
		),
	)
