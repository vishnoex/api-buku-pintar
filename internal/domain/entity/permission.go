package entity

import "time"

// Permission represents a permission in the system for RBAC (Role-Based Access Control)
// Clean Architecture: Entity layer, no dependencies on infrastructure
type Permission struct {
	ID          string    `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Resource    string    `db:"resource" json:"resource"`
	Action      string    `db:"action" json:"action"`
	Description *string   `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

// ResourceType represents the type of resource in the system
type ResourceType string

const (
	ResourceUser        ResourceType = "user"
	ResourceRole        ResourceType = "role"
	ResourcePermission  ResourceType = "permission"
	ResourceCategory    ResourceType = "category"
	ResourceBanner      ResourceType = "banner"
	ResourceEbook       ResourceType = "ebook"
	ResourceSummary     ResourceType = "summary"
	ResourceArticle     ResourceType = "article"
	ResourceInspiration ResourceType = "inspiration"
	ResourceAuthor      ResourceType = "author"
	ResourcePayment     ResourceType = "payment"
	ResourceComment     ResourceType = "comment"
	ResourceSEO         ResourceType = "seo"
)

// ActionType represents the type of action that can be performed on a resource
type ActionType string

const (
	ActionCreate ActionType = "create"
	ActionRead   ActionType = "read"
	ActionUpdate ActionType = "update"
	ActionDelete ActionType = "delete"
	ActionList   ActionType = "list"
	ActionManage ActionType = "manage" // Full access to the resource
)

// Common permission name patterns
const (
	// User permissions
	PermissionUserCreate = "user:create"
	PermissionUserRead   = "user:read"
	PermissionUserUpdate = "user:update"
	PermissionUserDelete = "user:delete"
	PermissionUserList   = "user:list"
	PermissionUserManage = "user:manage"

	// Role permissions
	PermissionRoleCreate = "role:create"
	PermissionRoleRead   = "role:read"
	PermissionRoleUpdate = "role:update"
	PermissionRoleDelete = "role:delete"
	PermissionRoleList   = "role:list"
	PermissionRoleManage = "role:manage"

	// Category permissions
	PermissionCategoryCreate = "category:create"
	PermissionCategoryRead   = "category:read"
	PermissionCategoryUpdate = "category:update"
	PermissionCategoryDelete = "category:delete"
	PermissionCategoryList   = "category:list"
	PermissionCategoryManage = "category:manage"

	// Banner permissions
	PermissionBannerCreate = "banner:create"
	PermissionBannerRead   = "banner:read"
	PermissionBannerUpdate = "banner:update"
	PermissionBannerDelete = "banner:delete"
	PermissionBannerList   = "banner:list"
	PermissionBannerManage = "banner:manage"

	// Ebook permissions
	PermissionEbookCreate = "ebook:create"
	PermissionEbookRead   = "ebook:read"
	PermissionEbookUpdate = "ebook:update"
	PermissionEbookDelete = "ebook:delete"
	PermissionEbookList   = "ebook:list"
	PermissionEbookManage = "ebook:manage"

	// Summary permissions
	PermissionSummaryCreate = "summary:create"
	PermissionSummaryRead   = "summary:read"
	PermissionSummaryUpdate = "summary:update"
	PermissionSummaryDelete = "summary:delete"
	PermissionSummaryList   = "summary:list"
	PermissionSummaryManage = "summary:manage"

	// Article permissions
	PermissionArticleCreate = "article:create"
	PermissionArticleRead   = "article:read"
	PermissionArticleUpdate = "article:update"
	PermissionArticleDelete = "article:delete"
	PermissionArticleList   = "article:list"
	PermissionArticleManage = "article:manage"

	// Inspiration permissions
	PermissionInspirationCreate = "inspiration:create"
	PermissionInspirationRead   = "inspiration:read"
	PermissionInspirationUpdate = "inspiration:update"
	PermissionInspirationDelete = "inspiration:delete"
	PermissionInspirationList   = "inspiration:list"
	PermissionInspirationManage = "inspiration:manage"

	// Author permissions
	PermissionAuthorCreate = "author:create"
	PermissionAuthorRead   = "author:read"
	PermissionAuthorUpdate = "author:update"
	PermissionAuthorDelete = "author:delete"
	PermissionAuthorList   = "author:list"
	PermissionAuthorManage = "author:manage"

	// Payment permissions
	PermissionPaymentCreate = "payment:create"
	PermissionPaymentRead   = "payment:read"
	PermissionPaymentUpdate = "payment:update"
	PermissionPaymentDelete = "payment:delete"
	PermissionPaymentList   = "payment:list"
	PermissionPaymentManage = "payment:manage"

	// Permission permissions (meta-permissions for managing permissions)
	PermissionPermissionCreate = "permission:create"
	PermissionPermissionRead   = "permission:read"
	PermissionPermissionUpdate = "permission:update"
	PermissionPermissionDelete = "permission:delete"
	PermissionPermissionList   = "permission:list"
	PermissionPermissionManage = "permission:manage"

	// Comment permissions
	PermissionCommentCreate = "comment:create"
	PermissionCommentRead   = "comment:read"
	PermissionCommentUpdate = "comment:update"
	PermissionCommentDelete = "comment:delete"
	PermissionCommentList   = "comment:list"
	PermissionCommentManage = "comment:manage"

	// SEO permissions
	PermissionSEOCreate = "seo:create"
	PermissionSEORead   = "seo:read"
	PermissionSEOUpdate = "seo:update"
	PermissionSEODelete = "seo:delete"
	PermissionSEOList   = "seo:list"
	PermissionSEOManage = "seo:manage"
)
