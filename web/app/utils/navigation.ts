import type {
  NavigationCatalog,
  NavigationItemKind,
  NavigationResult,
} from "../types/navigation";

const kindLabels: Record<NavigationItemKind, string> = {
  official: "官方",
  tool: "工具",
  community: "社区",
  learning: "学习",
  resource: "资源",
  reference: "参考",
  research: "研究",
};

function normalize(value: string) {
  return value.normalize("NFKC").trim().toLocaleLowerCase("zh-CN");
}

export function domainFromUrl(url: string) {
  return new URL(url).hostname.replace(/^www\./, "");
}

export function kindLabel(kind: NavigationItemKind) {
  return kindLabels[kind];
}

export function flattenNavigation(
  catalog: NavigationCatalog,
): NavigationResult[] {
  return catalog.categories.flatMap((category) =>
    category.groups.flatMap((group) =>
      group.items.map((item) => ({
        item,
        categoryId: category.id,
        categoryTitle: category.title,
        categoryIcon: category.icon,
        groupId: group.id,
        groupTitle: group.title,
        domain: domainFromUrl(item.url),
        searchText: normalize(
          [
            item.title,
            item.description,
            item.url,
            item.tags.join(" "),
            item.keywords?.join(" ") ?? "",
            category.title,
            group.title,
          ].join(" "),
        ),
      })),
    ),
  );
}

export function searchNavigation(entries: NavigationResult[], query: string) {
  const terms = normalize(query).split(/\s+/).filter(Boolean);
  if (!terms.length) return [];

  return entries
    .filter((entry) => terms.every((term) => entry.searchText.includes(term)))
    .sort((left, right) => {
      const normalizedQuery = normalize(query);
      const leftTitle = normalize(left.item.title);
      const rightTitle = normalize(right.item.title);
      const leftScore =
        leftTitle === normalizedQuery
          ? 0
          : leftTitle.startsWith(normalizedQuery)
            ? 1
            : left.item.featured
              ? 2
              : 3;
      const rightScore =
        rightTitle === normalizedQuery
          ? 0
          : rightTitle.startsWith(normalizedQuery)
            ? 1
            : right.item.featured
              ? 2
              : 3;
      return (
        leftScore - rightScore ||
        left.item.title.localeCompare(right.item.title, "zh-CN")
      );
    });
}
