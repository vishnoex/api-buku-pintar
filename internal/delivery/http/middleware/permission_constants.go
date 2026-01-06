package middleware

import "buku-pintar/internal/domain/entity"

// Permission constants are imported from entity package to avoid duplication
// Entity package is the single source of truth for permission constants
// This file only contains middleware-specific additions like:
// - Permission groups (Admin, Editor, Reader, etc.)
// - Permission descriptions
// - Helper functions

// Common permission groups for convenience
var (
	// AdminPermissions - All permissions an admin should have
	AdminPermissions = []string{
		entity.PermissionEbookManage,
		entity.PermissionArticleManage,
		entity.PermissionCategoryManage,
		entity.PermissionBannerManage,
		entity.PermissionUserManage,
		entity.PermissionRoleManage,
		entity.PermissionPermissionManage,
		entity.PermissionPaymentManage,
		entity.PermissionCommentManage,
		entity.PermissionSEOManage,
		entity.PermissionSummaryManage,
	}

	// EditorPermissions - Permissions for content editors
	EditorPermissions = []string{
		entity.PermissionEbookCreate,
		entity.PermissionEbookRead,
		entity.PermissionEbookUpdate,
		entity.PermissionArticleCreate,
		entity.PermissionArticleRead,
		entity.PermissionArticleUpdate,
		entity.PermissionSummaryCreate,
		entity.PermissionSummaryRead,
		entity.PermissionSummaryUpdate,
		entity.PermissionCategoryRead,
		entity.PermissionSEOUpdate,
	}

	// ReaderPermissions - Basic read permissions
	ReaderPermissions = []string{
		entity.PermissionEbookRead,
		entity.PermissionArticleRead,
		entity.PermissionCategoryRead,
		entity.PermissionBannerRead,
		entity.PermissionSummaryRead,
	}

	// PremiumPermissions - Enhanced read permissions for premium users
	PremiumPermissions = []string{
		entity.PermissionEbookRead,
		entity.PermissionArticleRead,
		entity.PermissionCategoryRead,
		entity.PermissionBannerRead,
		entity.PermissionSummaryRead,
		entity.PermissionCommentCreate,
		entity.PermissionCommentRead,
		entity.PermissionCommentUpdate, // Own comments only
	}

	// ContentManagerPermissions - Permissions for managing content
	ContentManagerPermissions = []string{
		entity.PermissionEbookManage,
		entity.PermissionArticleManage,
		entity.PermissionSummaryManage,
		entity.PermissionCategoryManage,
		entity.PermissionSEOManage,
	}

	// UserManagerPermissions - Permissions for user management
	UserManagerPermissions = []string{
		entity.PermissionUserManage,
		entity.PermissionRoleRead,
		entity.PermissionPermissionRead,
	}

	// PaymentManagerPermissions - Permissions for payment operations
	PaymentManagerPermissions = []string{
		entity.PermissionPaymentManage,
		entity.PermissionUserRead,
	}
)

// PermissionDescription provides human-readable descriptions for permissions
var PermissionDescription = map[string]string{
	// Ebook - from entity package
	entity.PermissionEbookCreate: "Create new ebooks",
	entity.PermissionEbookRead:   "View ebooks",
	entity.PermissionEbookUpdate: "Edit existing ebooks",
	entity.PermissionEbookDelete: "Delete ebooks",
	entity.PermissionEbookList:   "List ebooks",
	entity.PermissionEbookManage: "Full ebook management",

	// Article - from entity package
	entity.PermissionArticleCreate: "Create new articles",
	entity.PermissionArticleRead:   "View articles",
	entity.PermissionArticleUpdate: "Edit existing articles",
	entity.PermissionArticleDelete: "Delete articles",
	entity.PermissionArticleList:   "List articles",
	entity.PermissionArticleManage: "Full article management",

	// Category - from entity package
	entity.PermissionCategoryCreate: "Create new categories",
	entity.PermissionCategoryRead:   "View categories",
	entity.PermissionCategoryUpdate: "Edit existing categories",
	entity.PermissionCategoryDelete: "Delete categories",
	entity.PermissionCategoryList:   "List categories",
	entity.PermissionCategoryManage: "Full category management",

	// Banner - from entity package
	entity.PermissionBannerCreate: "Create new banners",
	entity.PermissionBannerRead:   "View banners",
	entity.PermissionBannerUpdate: "Edit existing banners",
	entity.PermissionBannerDelete: "Delete banners",
	entity.PermissionBannerList:   "List banners",
	entity.PermissionBannerManage: "Full banner management",

	// User - from entity package
	entity.PermissionUserCreate: "Create new users",
	entity.PermissionUserRead:   "View user information",
	entity.PermissionUserUpdate: "Edit user information",
	entity.PermissionUserDelete: "Delete users",
	entity.PermissionUserList:   "List users",
	entity.PermissionUserManage: "Full user management",

	// Role - from entity package
	entity.PermissionRoleCreate: "Create new roles",
	entity.PermissionRoleRead:   "View roles",
	entity.PermissionRoleUpdate: "Edit existing roles",
	entity.PermissionRoleDelete: "Delete roles",
	entity.PermissionRoleList:   "List roles",
	entity.PermissionRoleManage: "Full role management",

	// Permission - from entity package
	entity.PermissionPermissionCreate: "Create new permissions",
	entity.PermissionPermissionRead:   "View permissions",
	entity.PermissionPermissionUpdate: "Edit existing permissions",
	entity.PermissionPermissionDelete: "Delete permissions",
	entity.PermissionPermissionList:   "List permissions",
	entity.PermissionPermissionManage: "Full permission management",

	// Payment - from entity package
	entity.PermissionPaymentCreate: "Process new payments",
	entity.PermissionPaymentRead:   "View payment information",
	entity.PermissionPaymentUpdate: "Update payment status",
	entity.PermissionPaymentDelete: "Delete payment records",
	entity.PermissionPaymentList:   "List payments",
	entity.PermissionPaymentManage: "Full payment management",

	// Comment - from entity package
	entity.PermissionCommentCreate: "Create comments",
	entity.PermissionCommentRead:   "View comments",
	entity.PermissionCommentUpdate: "Edit comments",
	entity.PermissionCommentDelete: "Delete comments",
	entity.PermissionCommentList:   "List comments",
	entity.PermissionCommentManage: "Full comment management",

	// SEO - from entity package
	entity.PermissionSEOCreate: "Create SEO metadata",
	entity.PermissionSEORead:   "View SEO metadata",
	entity.PermissionSEOUpdate: "Edit SEO metadata",
	entity.PermissionSEODelete: "Delete SEO metadata",
	entity.PermissionSEOList:   "List SEO metadata",
	entity.PermissionSEOManage: "Full SEO management",

	// Summary - from entity package
	entity.PermissionSummaryCreate: "Create ebook summaries",
	entity.PermissionSummaryRead:   "View ebook summaries",
	entity.PermissionSummaryUpdate: "Edit ebook summaries",
	entity.PermissionSummaryDelete: "Delete ebook summaries",
	entity.PermissionSummaryList:   "List ebook summaries",
	entity.PermissionSummaryManage: "Full summary management",

	// Inspiration - from entity package
	entity.PermissionInspirationCreate: "Create inspirations",
	entity.PermissionInspirationRead:   "View inspirations",
	entity.PermissionInspirationUpdate: "Edit inspirations",
	entity.PermissionInspirationDelete: "Delete inspirations",
	entity.PermissionInspirationList:   "List inspirations",
	entity.PermissionInspirationManage: "Full inspiration management",

	// Author - from entity package
	entity.PermissionAuthorCreate: "Create authors",
	entity.PermissionAuthorRead:   "View authors",
	entity.PermissionAuthorUpdate: "Edit authors",
	entity.PermissionAuthorDelete: "Delete authors",
	entity.PermissionAuthorList:   "List authors",
	entity.PermissionAuthorManage: "Full author management",
}

// GetPermissionDescription returns a human-readable description for a permission
func GetPermissionDescription(permission string) string {
	if desc, exists := PermissionDescription[permission]; exists {
		return desc
	}
	return "Unknown permission"
}

// IsValidPermission checks if a permission name is valid
func IsValidPermission(permission string) bool {
	_, exists := PermissionDescription[permission]
	return exists
}

// GetPermissionsByCategory groups permissions by resource type
func GetPermissionsByCategory() map[string][]string {
	return map[string][]string{
		"ebook": {
			entity.PermissionEbookCreate,
			entity.PermissionEbookRead,
			entity.PermissionEbookUpdate,
			entity.PermissionEbookDelete,
			entity.PermissionEbookList,
			entity.PermissionEbookManage,
		},
		"article": {
			entity.PermissionArticleCreate,
			entity.PermissionArticleRead,
			entity.PermissionArticleUpdate,
			entity.PermissionArticleDelete,
			entity.PermissionArticleList,
			entity.PermissionArticleManage,
		},
		"category": {
			entity.PermissionCategoryCreate,
			entity.PermissionCategoryRead,
			entity.PermissionCategoryUpdate,
			entity.PermissionCategoryDelete,
			entity.PermissionCategoryList,
			entity.PermissionCategoryManage,
		},
		"banner": {
			entity.PermissionBannerCreate,
			entity.PermissionBannerRead,
			entity.PermissionBannerUpdate,
			entity.PermissionBannerDelete,
			entity.PermissionBannerList,
			entity.PermissionBannerManage,
		},
		"user": {
			entity.PermissionUserCreate,
			entity.PermissionUserRead,
			entity.PermissionUserUpdate,
			entity.PermissionUserDelete,
			entity.PermissionUserList,
			entity.PermissionUserManage,
		},
		"role": {
			entity.PermissionRoleCreate,
			entity.PermissionRoleRead,
			entity.PermissionRoleUpdate,
			entity.PermissionRoleDelete,
			entity.PermissionRoleList,
			entity.PermissionRoleManage,
		},
		"permission": {
			entity.PermissionPermissionCreate,
			entity.PermissionPermissionRead,
			entity.PermissionPermissionUpdate,
			entity.PermissionPermissionDelete,
			entity.PermissionPermissionList,
			entity.PermissionPermissionManage,
		},
		"payment": {
			entity.PermissionPaymentCreate,
			entity.PermissionPaymentRead,
			entity.PermissionPaymentUpdate,
			entity.PermissionPaymentDelete,
			entity.PermissionPaymentList,
			entity.PermissionPaymentManage,
		},
		"comment": {
			entity.PermissionCommentCreate,
			entity.PermissionCommentRead,
			entity.PermissionCommentUpdate,
			entity.PermissionCommentDelete,
			entity.PermissionCommentList,
			entity.PermissionCommentManage,
		},
		"seo": {
			entity.PermissionSEOCreate,
			entity.PermissionSEORead,
			entity.PermissionSEOUpdate,
			entity.PermissionSEODelete,
			entity.PermissionSEOList,
			entity.PermissionSEOManage,
		},
		"summary": {
			entity.PermissionSummaryCreate,
			entity.PermissionSummaryRead,
			entity.PermissionSummaryUpdate,
			entity.PermissionSummaryDelete,
			entity.PermissionSummaryList,
			entity.PermissionSummaryManage,
		},
		"inspiration": {
			entity.PermissionInspirationCreate,
			entity.PermissionInspirationRead,
			entity.PermissionInspirationUpdate,
			entity.PermissionInspirationDelete,
			entity.PermissionInspirationList,
			entity.PermissionInspirationManage,
		},
		"author": {
			entity.PermissionAuthorCreate,
			entity.PermissionAuthorRead,
			entity.PermissionAuthorUpdate,
			entity.PermissionAuthorDelete,
			entity.PermissionAuthorList,
			entity.PermissionAuthorManage,
		},
	}
}
