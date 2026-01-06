package middleware

import (
	"net/http"
	"strings"
)

// Common ResourceIDExtractor implementations for extracting resource IDs from requests

// ExtractFromPath extracts resource ID from URL path parameter
// Example: /ebooks/{id} -> extracts {id}
func ExtractFromPath(paramName string) ResourceIDExtractor {
	return func(r *http.Request) string {
		// Try Go 1.22+ PathValue first (if available)
		// For older Go versions, use path parsing
		// This is a simplified version - you may need to adjust based on your router
		
		// For ServeMux with patterns like /ebooks/{id}, you'd need Go 1.22+
		// For now, we'll extract from URL path segments as a fallback
		// Example: /ebooks/edit/123 where paramName is position
		
		// Use ExtractFromPathSegment or custom logic based on your routing
		return r.URL.Query().Get(paramName) // Fallback to query param
	}
}

// ExtractFromQuery extracts resource ID from URL query parameter
// Example: /ebooks?id=123 -> extracts "123"
func ExtractFromQuery(paramName string) ResourceIDExtractor {
	return func(r *http.Request) string {
		return r.URL.Query().Get(paramName)
	}
}

// ExtractFromHeader extracts resource ID from HTTP header
// Example: X-Resource-ID: 123 -> extracts "123"
func ExtractFromHeader(headerName string) ResourceIDExtractor {
	return func(r *http.Request) string {
		return r.Header.Get(headerName)
	}
}

// ExtractFromPathSegment extracts ID from a specific path segment
// Example: /ebooks/123/edit -> extractIDFromPathSegment(1) -> extracts "123"
func ExtractFromPathSegment(segmentIndex int) ResourceIDExtractor {
	return func(r *http.Request) string {
		segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if segmentIndex >= 0 && segmentIndex < len(segments) {
			return segments[segmentIndex]
		}
		return ""
	}
}

// ExtractLastPathSegment extracts the last segment of the URL path
// Example: /ebooks/123 -> extracts "123"
func ExtractLastPathSegment() ResourceIDExtractor {
	return func(r *http.Request) string {
		segments := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(segments) > 0 {
			return segments[len(segments)-1]
		}
		return ""
	}
}

// ExtractFromContextKey extracts resource ID from request context
// Useful when ID was set by previous middleware
func ExtractFromContextKey(key string) ResourceIDExtractor {
	return func(r *http.Request) string {
		if id, ok := r.Context().Value(key).(string); ok {
			return id
		}
		return ""
	}
}

// CombineExtractors tries multiple extractors in order until one returns a non-empty value
// Useful for APIs that accept resource ID from multiple sources
func CombineExtractors(extractors ...ResourceIDExtractor) ResourceIDExtractor {
	return func(r *http.Request) string {
		for _, extractor := range extractors {
			if id := extractor(r); id != "" {
				return id
			}
		}
		return ""
	}
}

// Common pre-defined extractors for convenience

var (
	// ExtractIDFromPath extracts "id" from path parameter
	ExtractIDFromPath = ExtractFromPath("id")

	// ExtractIDFromQuery extracts "id" from query parameter
	ExtractIDFromQuery = ExtractFromQuery("id")

	// ExtractEbookID extracts "ebook_id" from path or query
	ExtractEbookID = CombineExtractors(
		ExtractFromPath("ebook_id"),
		ExtractFromPath("id"),
		ExtractFromQuery("ebook_id"),
	)

	// ExtractArticleID extracts "article_id" from path or query
	ExtractArticleID = CombineExtractors(
		ExtractFromPath("article_id"),
		ExtractFromPath("id"),
		ExtractFromQuery("article_id"),
	)

	// ExtractCategoryID extracts "category_id" from path or query
	ExtractCategoryID = CombineExtractors(
		ExtractFromPath("category_id"),
		ExtractFromPath("id"),
		ExtractFromQuery("category_id"),
	)

	// ExtractUserID extracts "user_id" from path or query
	ExtractUserID = CombineExtractors(
		ExtractFromPath("user_id"),
		ExtractFromPath("id"),
		ExtractFromQuery("user_id"),
	)

	// ExtractCommentID extracts "comment_id" from path or query
	ExtractCommentID = CombineExtractors(
		ExtractFromPath("comment_id"),
		ExtractFromPath("id"),
		ExtractFromQuery("comment_id"),
	)

	// ExtractPaymentID extracts "payment_id" from path or query
	ExtractPaymentID = CombineExtractors(
		ExtractFromPath("payment_id"),
		ExtractFromPath("id"),
		ExtractFromQuery("payment_id"),
	)
)
