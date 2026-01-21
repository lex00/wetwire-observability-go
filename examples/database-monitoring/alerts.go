package monitoring

import (
	"github.com/lex00/wetwire-observability-go/promql"
	"github.com/lex00/wetwire-observability-go/rules"
)

// PostgreSQL Alerts

// PGHighConnectionUsage fires when PostgreSQL connection usage exceeds 80%.
var PGHighConnectionUsage = rules.AlertingRule{
	Alert: "PostgreSQLHighConnectionUsage",
	Expr:  promql.GT(PGConnectionsUsedExpr, promql.Scalar(0.8)).String(),
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "database",
	},
	Annotations: map[string]string{
		"summary":     "PostgreSQL connection pool near capacity",
		"description": "Connection usage is {{ $value | humanizePercentage }} on {{ $labels.instance }}",
	},
}

// PGCriticalConnectionUsage fires when PostgreSQL connection usage exceeds 95%.
var PGCriticalConnectionUsage = rules.AlertingRule{
	Alert: "PostgreSQLCriticalConnectionUsage",
	Expr:  promql.GT(PGConnectionsUsedExpr, promql.Scalar(0.95)).String(),
	For:   2 * rules.Minute,
	Labels: map[string]string{
		"severity": "critical",
		"team":     "database",
	},
	Annotations: map[string]string{
		"summary":     "PostgreSQL connection pool exhausted",
		"description": "Connection usage is {{ $value | humanizePercentage }} on {{ $labels.instance }}",
		"runbook_url": "https://runbooks.example.com/postgres/connection-exhaustion",
	},
}

// PGReplicationLag fires when replication lag exceeds 30 seconds.
var PGReplicationLag = rules.AlertingRule{
	Alert: "PostgreSQLReplicationLag",
	Expr:  promql.GT(PGReplicationLagExpr, promql.Scalar(30)).String(),
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "database",
	},
	Annotations: map[string]string{
		"summary":     "PostgreSQL replication lag detected",
		"description": "Replica {{ $labels.instance }} is {{ $value | humanizeDuration }} behind primary",
	},
}

// PGDeadlocks fires when deadlocks are detected.
var PGDeadlocks = rules.AlertingRule{
	Alert: "PostgreSQLDeadlocks",
	Expr:  promql.GT(PGDeadlocksExpr, promql.Scalar(0)).String(),
	For:   1 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "database",
	},
	Annotations: map[string]string{
		"summary":     "PostgreSQL deadlocks detected",
		"description": "Database {{ $labels.datname }} has {{ $value }} deadlocks/sec",
	},
}

// PGLowCacheHitRatio fires when cache hit ratio drops below 90%.
var PGLowCacheHitRatio = rules.AlertingRule{
	Alert: "PostgreSQLLowCacheHitRatio",
	Expr:  promql.LT(PGCacheHitRatioExpr, promql.Scalar(0.9)).String(),
	For:   10 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "database",
	},
	Annotations: map[string]string{
		"summary":     "PostgreSQL cache hit ratio low",
		"description": "Cache hit ratio is {{ $value | humanizePercentage }} for {{ $labels.datname }}",
	},
}

// Redis Alerts

// RedisHighMemoryUsage fires when Redis memory usage exceeds 80%.
var RedisHighMemoryUsage = rules.AlertingRule{
	Alert: "RedisHighMemoryUsage",
	Expr:  promql.GT(RedisMemoryUsedExpr, promql.Scalar(0.8)).String(),
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "database",
	},
	Annotations: map[string]string{
		"summary":     "Redis memory usage high",
		"description": "Memory usage is {{ $value | humanizePercentage }} on {{ $labels.instance }}",
	},
}

// RedisCriticalMemoryUsage fires when Redis memory usage exceeds 95%.
var RedisCriticalMemoryUsage = rules.AlertingRule{
	Alert: "RedisCriticalMemoryUsage",
	Expr:  promql.GT(RedisMemoryUsedExpr, promql.Scalar(0.95)).String(),
	For:   2 * rules.Minute,
	Labels: map[string]string{
		"severity": "critical",
		"team":     "database",
	},
	Annotations: map[string]string{
		"summary":     "Redis memory nearly exhausted",
		"description": "Memory usage is {{ $value | humanizePercentage }} on {{ $labels.instance }}",
		"runbook_url": "https://runbooks.example.com/redis/memory-exhaustion",
	},
}

// RedisEvictions fires when Redis is evicting keys.
var RedisEvictions = rules.AlertingRule{
	Alert: "RedisEvictions",
	Expr:  promql.GT(RedisEvictionsExpr, promql.Scalar(0)).String(),
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "database",
	},
	Annotations: map[string]string{
		"summary":     "Redis evicting keys",
		"description": "Redis {{ $labels.instance }} is evicting {{ $value }} keys/sec",
	},
}

// RedisLowHitRatio fires when Redis hit ratio drops below 80%.
var RedisLowHitRatio = rules.AlertingRule{
	Alert: "RedisLowHitRatio",
	Expr:  promql.LT(RedisHitRatioExpr, promql.Scalar(0.8)).String(),
	For:   10 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "database",
	},
	Annotations: map[string]string{
		"summary":     "Redis cache hit ratio low",
		"description": "Hit ratio is {{ $value | humanizePercentage }} on {{ $labels.instance }}",
	},
}
