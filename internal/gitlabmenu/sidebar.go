package gitlabmenu

import (
	"encoding/json"
	"github.com/aborche/gitlab-menu-extender-demo/internal/conf"
	"github.com/aborche/gitlab-menu-extender-demo/internal/parsers"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
)

type GitlabSidebar struct {
	AdminMode                           AdminMode             `json:"admin_mode"`
	AdminURL                            string                `json:"admin_url"`
	AvatarURL                           string                `json:"avatar_url"`
	CanSignOut                          bool                  `json:"can_sign_out"`
	CanaryToggleCOMURL                  string                `json:"canary_toggle_com_url"`
	ContextSwitcherLinks                []ContextSwitcherLink `json:"context_switcher_links"`
	CreateNewMenuGroups                 []CreateNewMenuGroup  `json:"create_new_menu_groups"`
	CurrentContext                      CurrentContext        `json:"current_context"`
	CurrentContextHeader                string                `json:"current_context_header"`
	CurrentMenuItems                    []CurrentMenuItem     `json:"current_menu_items"`
	DisplayWhatsNew                     bool                  `json:"display_whats_new"`
	GitlabCOMAndCanary                  bool                  `json:"gitlab_com_and_canary"`
	GitlabCOMButNotCanary               bool                  `json:"gitlab_com_but_not_canary"`
	GitlabVersion                       GitlabVersion         `json:"gitlab_version"`
	GitlabVersionCheck                  interface{}           `json:"gitlab_version_check"`
	GroupsPath                          string                `json:"groups_path"`
	HasLinkToProfile                    bool                  `json:"has_link_to_profile"`
	IsAdmin                             bool                  `json:"is_admin"`
	IsImpersonating                     bool                  `json:"is_impersonating"`
	IsLoggedIn                          bool                  `json:"is_logged_in"`
	IssuesDashboardPath                 string                `json:"issues_dashboard_path"`
	LinkToProfile                       string                `json:"link_to_profile"`
	LogoURL                             interface{}           `json:"logo_url"`
	MergeRequestMenu                    []MergeRequestMenu    `json:"merge_request_menu"`
	Name                                string                `json:"name"`
	PanelType                           string                `json:"panel_type"`
	PinnedItems                         []interface{}         `json:"pinned_items"`
	ProjectsPath                        string                `json:"projects_path"`
	Search                              Search                `json:"search"`
	Settings                            Settings              `json:"settings"`
	ShortcutLinks                       []ShortcutLink        `json:"shortcut_links"`
	ShowVersionCheck                    bool                  `json:"show_version_check"`
	SignOutLink                         string                `json:"sign_out_link"`
	Status                              Status                `json:"status"`
	StopImpersonationPath               string                `json:"stop_impersonation_path"`
	SupportPath                         string                `json:"support_path"`
	TodosDashboardPath                  string                `json:"todos_dashboard_path"`
	TrackVisitsPath                     string                `json:"track_visits_path"`
	UpdatePinsURL                       string                `json:"update_pins_url"`
	UserCounts                          UserCounts            `json:"user_counts"`
	Username                            string                `json:"username"`
	WhatsNewMostRecentReleaseItemsCount int64                 `json:"whats_new_most_recent_release_items_count"`
	WhatsNewVersionDigest               string                `json:"whats_new_version_digest"`
}

type AdminMode struct {
	AdminModeActive         bool   `json:"admin_mode_active"`
	AdminModeFeatureEnabled bool   `json:"admin_mode_feature_enabled"`
	EnterAdminModeURL       string `json:"enter_admin_mode_url"`
	LeaveAdminModeURL       string `json:"leave_admin_mode_url"`
	UserIsAdmin             bool   `json:"user_is_admin"`
}

type ContextSwitcherLink struct {
	Icon  string `json:"icon"`
	Link  string `json:"link"`
	Title string `json:"title"`
}

type CreateNewMenuGroup struct {
	Items []CreateNewMenuGroupItem `json:"items"`
	Name  string                   `json:"name"`
}

type CreateNewMenuGroupItem struct {
	Component  interface{}      `json:"component"`
	ExtraAttrs PurpleExtraAttrs `json:"extraAttrs"`
	Href       string           `json:"href"`
	Text       string           `json:"text"`
}

type PurpleExtraAttrs struct {
	DataQACreateMenuItem string `json:"data-qa-create-menu-item"`
	DataTestid           string `json:"data-testid"`
	DataTrackAction      string `json:"data-track-action"`
	DataTrackLabel       string `json:"data-track-label"`
	DataTrackProperty    string `json:"data-track-property"`
}

type CurrentContext struct {
}

type CurrentMenuItem struct {
	ActiveRoutes *ActiveRoutes         `json:"active_routes,omitempty"`
	Avatar       interface{}           `json:"avatar"`
	EntityID     interface{}           `json:"entity_id"`
	Icon         string                `json:"icon"`
	ID           string                `json:"id"`
	Link         string                `json:"link"`
	LinkClasses  interface{}           `json:"link_classes"`
	PillCount    interface{}           `json:"pill_count"`
	Title        string                `json:"title"`
	AvatarShape  string                `json:"avatar_shape,omitempty"`
	IsActive     bool                  `json:"is_active,omitempty"`
	Items        []CurrentMenuItemItem `json:"items,omitempty"`
	Separated    bool                  `json:"separated,omitempty"`
}

type ActiveRoutes struct {
	Controller string `json:"controller"`
}

type CurrentMenuItemItem struct {
	Avatar      interface{} `json:"avatar"`
	AvatarShape string      `json:"avatar_shape,omitempty"`
	EntityID    interface{} `json:"entity_id"`
	Icon        interface{} `json:"icon"`
	ID          string      `json:"id"`
	IsActive    bool        `json:"is_active"`
	Link        string      `json:"link"`
	LinkClasses interface{} `json:"link_classes"`
	PillCount   *int        `json:"pill_count"`
	Title       string      `json:"title"`
}

type GitlabVersion struct {
	Major   int64  `json:"major"`
	Minor   int64  `json:"minor"`
	Patch   int64  `json:"patch"`
	SuffixS string `json:"suffix_s"`
}

type MergeRequestMenu struct {
	Items []MergeRequestMenuItem `json:"items"`
	Name  string                 `json:"name"`
}

type MergeRequestMenuItem struct {
	Count      int64            `json:"count"`
	ExtraAttrs FluffyExtraAttrs `json:"extraAttrs"`
	Href       string           `json:"href"`
	Text       string           `json:"text"`
	UserCount  string           `json:"userCount"`
}

type FluffyExtraAttrs struct {
	Class             string `json:"class"`
	DataTrackAction   string `json:"data-track-action"`
	DataTrackLabel    string `json:"data-track-label"`
	DataTrackProperty string `json:"data-track-property"`
}

type Search struct {
	AutocompletePath string        `json:"autocomplete_path"`
	IssuesPath       string        `json:"issues_path"`
	MrPath           string        `json:"mr_path"`
	SearchContext    SearchContext `json:"search_context"`
	SearchPath       string        `json:"search_path"`
}

type SearchContext struct {
	ForSnippets interface{} `json:"for_snippets"`
}

type Settings struct {
	HasSettings            bool   `json:"has_settings"`
	ProfilePath            string `json:"profile_path"`
	ProfilePreferencesPath string `json:"profile_preferences_path"`
}

type ShortcutLink struct {
	CSSClass string `json:"css_class"`
	Href     string `json:"href"`
	Title    string `json:"title"`
}

type Status struct {
	Availability string      `json:"availability"`
	Busy         interface{} `json:"busy"`
	CanUpdate    bool        `json:"can_update"`
	ClearAfter   interface{} `json:"clear_after"`
	Customized   interface{} `json:"customized"`
	Emoji        interface{} `json:"emoji"`
	Message      interface{} `json:"message"`
	MessageHTML  interface{} `json:"message_html"`
}

type UserCounts struct {
	AssignedIssues               int64 `json:"assigned_issues"`
	AssignedMergeRequests        int64 `json:"assigned_merge_requests"`
	LastUpdate                   int64 `json:"last_update"`
	ReviewRequestedMergeRequests int64 `json:"review_requested_merge_requests"`
	Todos                        int64 `json:"todos"`
}

type TabMenu struct {
	Name        string
	Branch      string
	Path        string
	ActiveStyle string
}

func TransformSidebar(c *gin.Context, cfg conf.Config, source *html.Node, fromFile bool) ([]byte, error) {
	// Get sidebar source from Gitlab page
	GitlabParsedSidebar, err := parsers.GetSideBarMenu(source)
	if err != nil {
		return nil, err
	}

	var SourceSideBar GitlabSidebar
	err = json.Unmarshal([]byte(GitlabParsedSidebar), &SourceSideBar)
	if err != nil {
		return nil, err
	}
	// Remove from menu unused items
	SourceSideBar.AdminMode.UserIsAdmin = false
	SourceSideBar.IsAdmin = false
	SourceSideBar.CanSignOut = false
	SourceSideBar.Status.CanUpdate = false
	SourceSideBar.HasLinkToProfile = false
	SourceSideBar.Settings.HasSettings = false
	SourceSideBar.ContextSwitcherLinks = []ContextSwitcherLink{}
	SourceSideBar.ShortcutLinks = []ShortcutLink{}
	SourceSideBar.MergeRequestMenu = []MergeRequestMenu{}
	SourceSideBar.CreateNewMenuGroups = []CreateNewMenuGroup{}
	SourceSideBar.DisplayWhatsNew = false
	if fromFile {
		SourceSideBar.CurrentMenuItems = BuildSidebarMenuFromFile(c, cfg)
	} else {
		SourceSideBar.CurrentMenuItems = BuildSidebarMenuFromStruct(c, cfg)
	}
	return json.Marshal(SourceSideBar)

}
