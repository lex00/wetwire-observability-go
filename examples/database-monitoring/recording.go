package monitoring

import "github.com/lex00/wetwire-observability-go/rules"

// PostgreSQL Recording Rules

// PGConnectionUsage5m pre-computes connection usage ratio.
var PGConnectionUsage5m = rules.RecordingRule{
	Record: "postgres:connection_usage:ratio",
	Expr:   PGConnectionsUsedExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
	},
}

// PGCacheHitRatio5m pre-computes cache hit ratio.
var PGCacheHitRatio5m = rules.RecordingRule{
	Record: "postgres:cache_hit_ratio:5m",
	Expr:   PGCacheHitRatioExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
	},
}

// PGDeadlockRate5m pre-computes deadlock rate.
var PGDeadlockRate5m = rules.RecordingRule{
	Record: "postgres:deadlocks:rate5m",
	Expr:   PGDeadlocksExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
	},
}

// PGReplicationLagMax pre-computes maximum replication lag.
var PGReplicationLagMax = rules.RecordingRule{
	Record: "postgres:replication_lag:max",
	Expr:   PGReplicationLagExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
	},
}

// Redis Recording Rules

// RedisMemoryUsage pre-computes memory usage ratio.
var RedisMemoryUsage = rules.RecordingRule{
	Record: "redis:memory_usage:ratio",
	Expr:   RedisMemoryUsedExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
	},
}

// RedisHitRatio5m pre-computes hit ratio.
var RedisHitRatio5m = rules.RecordingRule{
	Record: "redis:hit_ratio:5m",
	Expr:   RedisHitRatioExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
	},
}

// RedisCommandsRate5m pre-computes commands per second.
var RedisCommandsRate5m = rules.RecordingRule{
	Record: "redis:commands:rate5m",
	Expr:   RedisCommandsPerSecExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
	},
}

// RedisEvictionRate5m pre-computes eviction rate.
var RedisEvictionRate5m = rules.RecordingRule{
	Record: "redis:evictions:rate5m",
	Expr:   RedisEvictionsExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
	},
}
