package monitoring

import "github.com/lex00/wetwire-observability-go/alertmanager"

// InhibitionRules define which alerts should be muted when other alerts are firing.

// CriticalInhibitsWarning mutes warning alerts when a critical alert is firing
// for the same alertname and service.
var CriticalInhibitsWarning = alertmanager.NewInhibitRule().
	WithSourceMatchers(alertmanager.SeverityCritical()).
	WithTargetMatchers(alertmanager.SeverityWarning()).
	WithEqual("alertname", "service")

// CriticalInhibitsInfo mutes info alerts when a critical alert is firing.
var CriticalInhibitsInfo = alertmanager.NewInhibitRule().
	WithSourceMatchers(alertmanager.SeverityCritical()).
	WithTargetMatchers(alertmanager.SeverityInfo()).
	WithEqual("alertname", "service")

// WarningInhibitsInfo mutes info alerts when a warning alert is firing.
var WarningInhibitsInfo = alertmanager.NewInhibitRule().
	WithSourceMatchers(alertmanager.SeverityWarning()).
	WithTargetMatchers(alertmanager.SeverityInfo()).
	WithEqual("alertname", "service")

// ClusterDownInhibitsAll mutes all alerts from a cluster when the cluster is down.
var ClusterDownInhibitsAll = alertmanager.NewInhibitRule().
	WithSourceMatchers(alertmanager.Alertname("ClusterDown")).
	WithTargetMatchers(alertmanager.NotEq("alertname", "ClusterDown")).
	WithEqual("cluster")

// NodeDownInhibitsPodAlerts mutes pod alerts when the node is down.
var NodeDownInhibitsPodAlerts = alertmanager.NewInhibitRule().
	WithSourceMatchers(alertmanager.Alertname("NodeDown")).
	WithTargetMatchers(alertmanager.AlertnameRegex("Pod.*")).
	WithEqual("node")

// MaintenanceInhibitsAll mutes all alerts during maintenance.
var MaintenanceInhibitsAll = alertmanager.NewInhibitRule().
	WithSourceMatchers(alertmanager.Alertname("MaintenanceMode")).
	WithTargetMatchers(alertmanager.NotEq("alertname", "MaintenanceMode")).
	WithEqual("service")

// InhibitionRules is the list of all inhibit rules.
var InhibitionRules = []*alertmanager.InhibitRule{
	CriticalInhibitsWarning,
	CriticalInhibitsInfo,
	WarningInhibitsInfo,
	ClusterDownInhibitsAll,
	NodeDownInhibitsPodAlerts,
	MaintenanceInhibitsAll,
}
