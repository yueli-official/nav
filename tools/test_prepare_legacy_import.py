import importlib.util
from pathlib import Path
import unittest

spec = importlib.util.spec_from_file_location("legacy", Path(__file__).with_name("prepare-legacy-import.py"))
legacy = importlib.util.module_from_spec(spec)
spec.loader.exec_module(legacy)


class LegacyPreparationTest(unittest.TestCase):
    def rows(self):
        return [{"category": {"cid": cid, "title": cid}, "groups": [{"group": {"gid": "common", "title": "Group"}, "collections": [{"title": "O'Reilly", "link": "https://example.com/", "tags": ["Books"]}]}]} for cid in ["one", "two"]]

    def test_scopes_repeated_group_keys_and_preserves_duplicate_links(self):
        categories, groups, links, errors, duplicates = legacy.prepare(self.rows())
        self.assertEqual((len(categories), len(groups), len(links)), (2, 2, 2))
        self.assertEqual(errors, [])
        self.assertEqual(len({row[0] for row in groups}), 2)
        self.assertEqual(len({row[0] for row in links}), 2)
        self.assertEqual(len(duplicates["https://example.com/"]), 2)
        self.assertEqual(legacy.prepare(self.rows()), legacy.prepare(self.rows()))

    def test_rejects_unsafe_or_missing_links_without_silent_loss(self):
        rows = self.rows()
        rows[0]["groups"][0]["collections"][0]["link"] = "javascript:alert(1)"
        result = legacy.prepare(rows)
        self.assertEqual(len(result[2]), 2)
        self.assertEqual(len(result[3]), 1)

    def test_reports_category_and_group_collisions(self):
        rows = self.rows()
        rows.append(rows[0])
        errors = legacy.prepare(rows)[3]
        self.assertEqual(len(errors), 2)

    def test_sql_literal_preserves_quotes_and_blocks_nul(self):
        self.assertEqual(legacy.literal("O'Reilly"), "'O''Reilly'")
        with self.assertRaises(ValueError):
            legacy.literal("invalid\0text")


if __name__ == "__main__":
    unittest.main()
