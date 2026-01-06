package middleware

import (
	"fmt"
	"sync"
	"time"
)

// PermissionAuditEntry represents a single permission check audit log entry
type PermissionAuditEntry struct {
	Timestamp  time.Time
	UserID     string
	Permission string
	Granted    bool
	Reason     string
	Duration   time.Duration
	RequestID  string // Optional: for tracing requests
	IPAddress  string // Optional: for security audits
}

// PermissionAuditLogger handles logging of permission checks for audit trails
type PermissionAuditLogger struct {
	entries []PermissionAuditEntry
	mu      sync.RWMutex
	maxSize int
}

// NewPermissionAuditLogger creates a new audit logger
func NewPermissionAuditLogger() *PermissionAuditLogger {
	return &PermissionAuditLogger{
		entries: make([]PermissionAuditEntry, 0),
		maxSize: 10000, // Keep last 10,000 entries in memory
	}
}

// Log adds a new entry to the audit log
func (l *PermissionAuditLogger) Log(entry PermissionAuditEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.entries = append(l.entries, entry)

	// Trim if exceeds max size (keep most recent entries)
	if len(l.entries) > l.maxSize {
		l.entries = l.entries[len(l.entries)-l.maxSize:]
	}
}

// GetEntries returns all audit log entries
func (l *PermissionAuditLogger) GetEntries() []PermissionAuditEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// Return a copy to prevent external modification
	entries := make([]PermissionAuditEntry, len(l.entries))
	copy(entries, l.entries)
	return entries
}

// GetEntriesByUser returns audit log entries for a specific user
func (l *PermissionAuditLogger) GetEntriesByUser(userID string) []PermissionAuditEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var userEntries []PermissionAuditEntry
	for _, entry := range l.entries {
		if entry.UserID == userID {
			userEntries = append(userEntries, entry)
		}
	}
	return userEntries
}

// GetDeniedEntries returns all denied permission attempts
func (l *PermissionAuditLogger) GetDeniedEntries() []PermissionAuditEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var denied []PermissionAuditEntry
	for _, entry := range l.entries {
		if !entry.Granted {
			denied = append(denied, entry)
		}
	}
	return denied
}

// GetEntriesSince returns entries since a specific timestamp
func (l *PermissionAuditLogger) GetEntriesSince(since time.Time) []PermissionAuditEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var recent []PermissionAuditEntry
	for _, entry := range l.entries {
		if entry.Timestamp.After(since) {
			recent = append(recent, entry)
		}
	}
	return recent
}

// GetStats returns statistics about permission checks
func (l *PermissionAuditLogger) GetStats() PermissionAuditStats {
	l.mu.RLock()
	defer l.mu.RUnlock()

	stats := PermissionAuditStats{
		TotalChecks: len(l.entries),
		UniqueUsers: make(map[string]bool),
		ByPermission: make(map[string]int),
	}

	for _, entry := range l.entries {
		if entry.Granted {
			stats.GrantedCount++
		} else {
			stats.DeniedCount++
		}

		stats.UniqueUsers[entry.UserID] = true
		stats.ByPermission[entry.Permission]++
		stats.TotalDuration += entry.Duration
	}

	stats.UniqueUserCount = len(stats.UniqueUsers)
	if stats.TotalChecks > 0 {
		stats.AverageDuration = stats.TotalDuration / time.Duration(stats.TotalChecks)
	}

	return stats
}

// Clear removes all entries from the audit log
func (l *PermissionAuditLogger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.entries = make([]PermissionAuditEntry, 0)
}

// Export returns all entries in a format suitable for external storage
func (l *PermissionAuditLogger) Export() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	exported := make([]string, len(l.entries))
	for i, entry := range l.entries {
		status := "GRANTED"
		if !entry.Granted {
			status = "DENIED"
		}
		exported[i] = fmt.Sprintf("%s | User: %s | Permission: %s | Status: %s | Reason: %s | Duration: %v",
			entry.Timestamp.Format(time.RFC3339),
			entry.UserID,
			entry.Permission,
			status,
			entry.Reason,
			entry.Duration,
		)
	}
	return exported
}

// PermissionAuditStats holds statistics about permission checks
type PermissionAuditStats struct {
	TotalChecks     int
	GrantedCount    int
	DeniedCount     int
	UniqueUserCount int
	UniqueUsers     map[string]bool
	ByPermission    map[string]int
	TotalDuration   time.Duration
	AverageDuration time.Duration
}

// String returns a human-readable representation of the stats
func (s PermissionAuditStats) String() string {
	return fmt.Sprintf(
		"Permission Check Stats:\n"+
			"  Total Checks: %d\n"+
			"  Granted: %d (%.1f%%)\n"+
			"  Denied: %d (%.1f%%)\n"+
			"  Unique Users: %d\n"+
			"  Average Duration: %v",
		s.TotalChecks,
		s.GrantedCount, float64(s.GrantedCount)/float64(s.TotalChecks)*100,
		s.DeniedCount, float64(s.DeniedCount)/float64(s.TotalChecks)*100,
		s.UniqueUserCount,
		s.AverageDuration,
	)
}
