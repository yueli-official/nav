"""Prepare a reviewable legacy JSON import; never connects to a database."""
import argparse
import hashlib
import json
from pathlib import Path
from urllib.parse import urlsplit


def prepare(rows):
    categories, groups, links, errors = [], [], [], []
    seen_categories, seen_groups = set(), set()
    urls = {}
    for category_index, row in enumerate(rows):
        category = row["category"]
        cid = "legacy-" + str(category["cid"])
        if cid in seen_categories:
            errors.append(f"Duplicate category: {cid}")
        seen_categories.add(cid)
        categories.append((cid, category["title"], "", "i-tabler-folder", category.get("order", category_index)))
        for group_index, entry in enumerate(row.get("groups", [])):
            group = entry["group"]
            gid = cid + "-" + str(group["gid"])
            if gid in seen_groups:
                errors.append(f"Duplicate scoped group: {gid}")
            seen_groups.add(gid)
            groups.append((gid, cid, group["title"], "", group.get("order", group_index)))
            for index, link in enumerate(entry.get("collections", [])):
                url = str(link.get("link", "")).strip()
                parsed = urlsplit(url)
                title = str(link.get("title", "")).strip()
                if not title or parsed.scheme not in ("http", "https") or not parsed.netloc:
                    errors.append(f"Invalid link at {cid}/{gid}/{index}: {title}")
                digest = hashlib.sha256(f"{cid}\0{gid}\0{index}\0{url}".encode()).hexdigest()[:24]
                tags = link.get("tags") or []
                if not isinstance(tags, list) or not all(isinstance(tag, str) for tag in tags):
                    errors.append(f"Invalid tags at {gid}/{index}")
                links.append(("legacy-" + digest, cid, gid, title, url, str(link.get("description") or ""), json.dumps(tags, ensure_ascii=False), "[]", "reference", False, "published", index))
                urls.setdefault(url, []).append(title)
    return categories, groups, links, errors, {url: titles for url, titles in urls.items() if len(titles) > 1}


def literal(value):
    if isinstance(value, bool):
        return "TRUE" if value else "FALSE"
    if isinstance(value, int):
        return str(value)
    if "\0" in value:
        raise ValueError("NUL is not valid PostgreSQL text")
    return "'" + value.replace("'", "''") + "'"


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--input", type=Path, required=True)
    parser.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()
    raw = args.input.read_bytes()
    data = json.loads(raw.decode("utf-8-sig"))
    rows = data["nav"] if isinstance(data, dict) else data
    categories, groups, links, errors, duplicates = prepare(rows)
    args.output.mkdir(parents=True, exist_ok=True)
    report = {"sourceSha256": hashlib.sha256(raw).hexdigest(), "categories": len(categories), "groups": len(groups), "links": len(links), "errors": errors, "duplicateURLs": duplicates, "writesDatabase": False, "behavior": "Preserve all legacy links and order; scoped group IDs; deterministic link IDs. No deletions or updates of existing rows. Existing seeded links require separate overlap review."}
    (args.output / "report.json").write_text(json.dumps(report, ensure_ascii=False, indent=2), encoding="utf-8")
    if errors:
        raise SystemExit("Invalid legacy rows; inspect report.json. No SQL produced.")
    sql = ["\\set ON_ERROR_STOP on", "BEGIN;", "SELECT pg_advisory_xact_lock(hashtext('nav-legacy-import'));", "SET LOCAL standard_conforming_strings = on;"]
    tables = [("nav_categories", "id,title,description,icon,sort_order", categories), ("nav_groups", "id,category_id,title,description,sort_order", groups), ("nav_links", "id,category_id,group_id,title,url,description,tags,keywords,kind,featured,status,sort_order", links)]
    for table, columns, values in tables:
        for row in values:
            sql.append(f"INSERT INTO {table} ({columns}) VALUES ({','.join(literal(value) for value in row)}) ON CONFLICT (id) DO NOTHING;")
    sql.append("COMMIT;")
    (args.output / "import.sql").write_text("\n".join(sql) + "\n", encoding="utf-8")
    print(json.dumps({key: report[key] for key in ["sourceSha256", "categories", "groups", "links", "writesDatabase"]}))


if __name__ == "__main__":
    main()
