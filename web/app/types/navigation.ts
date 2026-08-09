import type {
  SiteProfileFormSchema,
  SiteProfileSnapshot,
} from "@yueli/site-profile/types";

export type NavigationItemKind =
  | "official"
  | "tool"
  | "community"
  | "learning"
  | "resource"
  | "reference"
  | "research";

export interface NavigationItem {
  id: string;
  title: string;
  url: string;
  description: string;
  tags: string[];
  keywords?: string[];
  kind: NavigationItemKind;
  featured: boolean;
  clickCount: number;
  faviconRevision?: string;
  lastClickedAt?: string;
}

export interface NavigationGroup {
  id: string;
  categoryId?: string;
  title: string;
  description: string;
  sortOrder: number;
  linkCount: number;
  items: NavigationItem[];
}

export interface NavigationCategory {
  id: string;
  title: string;
  description: string;
  icon: string;
  sortOrder: number;
  groups: NavigationGroup[];
}

export interface NavigationSiteCopy {
  revision: number;
  runtimeRevision: number;
  etag: string;
  name: string;
  title: string;
  description: string;
  searchPlaceholder: string;
  footerTagline: string;
}

export interface NavigationCatalog {
  version: 1;
  site: NavigationSiteCopy;
  categories: NavigationCategory[];
}

export interface NavigationStats {
  categoryCount: number;
  groupCount: number;
  linkCount: number;
}

export interface NavigationResponse extends NavigationCatalog {
  stats: NavigationStats;
}

export interface NavigationResult {
  item: NavigationItem;
  categoryId: string;
  categoryTitle: string;
  categoryIcon: string;
  groupId: string;
  groupTitle: string;
  domain: string;
  searchText: string;
}

export type NavigationStatus = "published" | "draft" | "archived";

export interface AdminNavigationLink extends NavigationItem {
  categoryId: string;
  groupId: string;
  status: NavigationStatus;
  sortOrder: number;
  publishedAt?: string;
  createdAt?: string;
  updatedAt?: string;
  healthStatus?: NavigationHealthStatus;
  lastCheckedAt?: string;
  healthHttpStatus?: number;
  healthLatencyMs?: number;
  healthError?: string;
  healthCheckExempt?: boolean;
}

export type NavigationHealthStatus =
  "unchecked" | "healthy" | "redirected" | "broken" | "timeout" | "error";

export interface NavigationGroupResponse {
  site: NavigationSiteCopy;
  category: NavigationCategory;
  group: NavigationGroup;
  items: NavigationItem[];
  total: number;
  page: number;
  size: number;
}

export interface NavigationHealthCounts {
  all: number;
  unchecked: number;
  healthy: number;
  redirected: number;
  broken: number;
  timeout: number;
  error: number;
  exempt: number;
}

export interface NavigationChecksResponse {
  links: AdminNavigationLink[];
  counts: NavigationHealthCounts;
  total: number;
  checkableTotal: number;
  page: number;
  size: number;
}

export type NavigationCheckJobStatus = "running" | "completed" | "failed";

export interface NavigationCheckJob {
  id: string;
  scope: "filtered" | "selected";
  status: NavigationCheckJobStatus;
  total: number;
  completed: number;
  startedAt: string;
  finishedAt?: string;
  error?: string;
}

export interface NavigationStartCheckJobResponse {
  job: NavigationCheckJob;
  reused: boolean;
}

export interface NavigationCheckJobResponse {
  job: NavigationCheckJob;
}

export interface NavigationCheckExemptionResponse {
  link: AdminNavigationLink;
}

export interface AdminNavigationResponse {
  links: AdminNavigationLink[];
  categories: NavigationCategory[];
  tags: NavigationTag[];
  counts: NavigationLifecycleCounts;
  total: number;
  page: number;
  size: number;
}

export interface NavigationTag {
  name: string;
  linkCount: number;
}

export interface NavigationLifecycleCounts {
  all: number;
  published: number;
  draft: number;
  archived: number;
}

export interface NavigationStructureResponse {
  categories: NavigationCategory[];
}

export interface NavigationTagsResponse {
  tags: NavigationTag[];
}

export interface NavigationSettingsResponse {
  settings: {
    snapshot: SiteProfileSnapshot;
    schema: SiteProfileFormSchema;
    searchPlaceholder: string;
    runtimeRevision: number;
    etag: string;
  };
}

export interface NavMeView {
  sub: string;
  userKey?: string;
  authenticated: boolean;
  isAdministrator: boolean;
  capabilities: string[];
  membership?: {
    status: "active" | "suspended";
    joinedAt: string;
    lastSeenAt: string;
  };
}

export interface NavigationMemberRole {
  key: string;
  displayName: string;
  source: string;
}

export interface NavigationMember {
  userKey: string;
  status: "active" | "suspended";
  displayName: string;
  handle: string;
  avatarMediaKey?: string;
  joinedAt: string;
  lastSeenAt: string;
  suspendedAt?: string;
  suspendedBy?: string;
  suspensionReason?: string;
  submissionCount: number;
  pendingApplications: number;
  roles: NavigationMemberRole[];
}

export interface NavigationMemberCounts {
  all: number;
  active: number;
  suspended: number;
}

export interface NavigationMemberRoleOption {
  key: string;
  displayName: string;
  status: string;
}

export interface NavigationMembersResponse {
  members: NavigationMember[];
  counts: NavigationMemberCounts;
  roles: NavigationMemberRoleOption[];
  total: number;
  page: number;
  size: number;
}
