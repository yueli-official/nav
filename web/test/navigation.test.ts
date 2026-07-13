import { describe, expect, test } from "vitest";
import {
  domainFromUrl,
  flattenNavigation,
  kindLabel,
  searchNavigation,
} from "../app/utils/navigation";
import type { NavigationCatalog } from "../app/types/navigation";

const catalog: NavigationCatalog = {
  version: 1,
  site: {
    name: "Test Nav",
    title: "Test",
    description: "Test catalog",
    searchPlaceholder: "Search",
  },
  categories: [
    {
      id: "develop",
      title: "开发工程",
      description: "开发资料",
      icon: "i-tabler-code",
      groups: [
        {
          id: "references",
          title: "文档",
          description: "权威文档",
          items: [
            {
              id: "mdn",
              title: "MDN Web Docs",
              url: "https://developer.mozilla.org/",
              description: "Web 平台文档",
              tags: ["JavaScript", "CSS"],
              kind: "reference",
              featured: true,
            },
          ],
        },
      ],
    },
  ],
};

describe("navigation utilities", () => {
  test("flattens hierarchy with searchable context", () => {
    const [entry] = flattenNavigation(catalog);
    expect(entry).toMatchObject({
      categoryTitle: "开发工程",
      groupTitle: "文档",
      domain: "developer.mozilla.org",
    });
    expect(entry?.searchText).toContain("javascript");
  });

  test("searches across title, tags and hierarchy", () => {
    const entries = flattenNavigation(catalog);
    expect(searchNavigation(entries, "MDN")).toHaveLength(1);
    expect(searchNavigation(entries, "开发 CSS")).toHaveLength(1);
    expect(searchNavigation(entries, "三维")).toHaveLength(0);
  });

  test("normalizes domains and kind labels", () => {
    expect(domainFromUrl("https://www.example.com/path")).toBe("example.com");
    expect(kindLabel("reference")).toBe("参考");
  });
});
