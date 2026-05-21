"""
check_routes.py
Compares routes defined in a Go Fiber routes.go file with those documented
in a swagger.yaml file, reporting any mismatches in either direction.

Handles:
  - Arbitrary nesting depth of .Group() calls
  - Groups with middleware args, e.g. v1.Group("", middleware1, middleware2)
  - Empty-string group prefixes (used purely to attach middleware)
  - Path params: converts Go :param  →  Swagger {param}
  - The /api/v1 base prefix from the v1 group
  - Inline block scoping with { … }  (does not affect extraction)
"""

import re
import sys
import yaml
from pathlib import Path
from collections import defaultdict


# ---------------------------------------------------------------------------
# Go route extraction
# ---------------------------------------------------------------------------

# var := someReceiver.Group("/prefix")          – captures (var, receiver, prefix)
# Also matches when the first arg is "" (empty-string group for middleware only)
GROUP_RE = re.compile(
    r'(\w+)\s*:=\s*(\w+)\.Group\(\s*"([^"]*)"'
)

# someReceiver.Get("/path", handler...)         – captures (receiver, METHOD, path)
ROUTE_RE = re.compile(
    r'(\w+)\.(Get|Post|Put|Delete|Patch|Head|Options)\(\s*"(/[^"]*)"'
)


def _normalise_go_path(path: str) -> str:
    """Convert Go Fiber :param style to OpenAPI {param} style, and strip trailing slashes."""
    path = re.sub(r":(\w+)", r"{\1}", path)
    # Strip trailing slash unless the path is literally "/"
    if path != "/" and path.endswith("/"):
        path = path.rstrip("/")
    return path


def extract_go_routes(file_path: Path) -> set[str]:
    """
    Parse routes.go and return the full set of normalised route paths.
    Works by building a prefix map for every group variable, then
    prepending the accumulated prefix whenever a route is registered.
    """
    source = file_path.read_text(encoding="utf-8")
    lines = source.splitlines()

    # variable_name -> accumulated prefix string
    # Seed with the two well-known top-level receivers
    group_prefix: dict[str, str] = {
        "app": "",   # app.Get(…) → no prefix beyond what groups add
        "v1": "/api/v1",
    }

    routes: set[str] = set()

    for line in lines:
        # ---- group declarations ----
        gm = GROUP_RE.search(line)
        if gm:
            var, parent, subprefix = gm.groups()
            parent_prefix = group_prefix.get(parent, "")
            group_prefix[var] = parent_prefix + subprefix
            continue  # a group line won't also be a route line

        # ---- route registrations ----
        rm = ROUTE_RE.search(line)
        if rm:
            receiver, _method, path = rm.groups()
            prefix = group_prefix.get(receiver, "")
            full = _normalise_go_path(prefix + path)
            routes.add(full)

    return routes


# ---------------------------------------------------------------------------
# Swagger route extraction
# ---------------------------------------------------------------------------

def extract_swagger_routes(file_path: Path) -> set[str]:
    """Return the set of paths defined in the swagger.yaml `paths` section."""
    with file_path.open(encoding="utf-8") as fh:
        spec = yaml.safe_load(fh)

    raw_paths: dict = spec.get("paths", {})
    routes: set[str] = set()

    for path, path_item in raw_paths.items():
        if not isinstance(path_item, dict):
            routes.add(path)
            continue

        # Flatten any incorrectly nested paths (e.g. the /admin/users/{id} block
        # that appears *inside* the /admin/users entry in this swagger file)
        for key, value in path_item.items():
            if key.startswith("/"):
                # This is a mis-indented nested path — treat it as a top-level path
                if isinstance(value, dict):
                    routes.add(key)
            
        routes.add(path)

    return routes


# ---------------------------------------------------------------------------
# Comparison & reporting
# ---------------------------------------------------------------------------

def _strip_api_v1(path: str) -> str:
    """Remove the /api/v1 prefix so we can compare against bare swagger paths."""
    if path.startswith("/api/v1"):
        return path[len("/api/v1"):]
    return path


def compare_and_report(go_file: Path, swagger_file: Path) -> int:
    """
    Main comparison logic.
    Returns 0 if routes are in sync, 1 if there are differences.
    """
    raw_go = extract_go_routes(go_file)
    swagger = extract_swagger_routes(swagger_file)

    # Swagger paths don't carry the /api/v1 prefix, so strip it from Go paths
    # before comparing.  Keep originals for display purposes.
    go_stripped: dict[str, str] = {_strip_api_v1(p): p for p in raw_go}

    only_in_go = {
        go_stripped[s]: s
        for s in go_stripped
        if s not in swagger
    }
    only_in_swagger = swagger - set(go_stripped)

    # ---- summary -------------------------------------------------------
    total_go = len(go_stripped)
    total_sw = len(swagger)
    matched = total_go - len(only_in_go)

    print("=" * 60)
    print("  ROUTE SYNC CHECK")
    print("=" * 60)
    print(f"  Go routes found    : {total_go}")
    print(f"  Swagger paths found: {total_sw}")
    print(f"  Matched            : {matched}")
    print(f"  Only in Go         : {len(only_in_go)}")
    print(f"  Only in Swagger    : {len(only_in_swagger)}")
    print("=" * 60)

    if only_in_go:
        print("\n🔴  Registered in Go but MISSING from Swagger:")
        for stripped, original in sorted(only_in_go.items()):
            print(f"    {original}  →  (swagger path would be: {stripped})")

    if only_in_swagger:
        print("\n🟡  Documented in Swagger but NOT registered in Go:")
        for path in sorted(only_in_swagger):
            print(f"    {path}")

    if not only_in_go and not only_in_swagger:
        print("\n✅  All routes are in sync!")
        return 0

    print()
    return 1


# ---------------------------------------------------------------------------
# Entry point
# ---------------------------------------------------------------------------

if __name__ == "__main__":
    import argparse

    parser = argparse.ArgumentParser(
        description="Check that Go Fiber routes and swagger.yaml are in sync."
    )
    parser.add_argument(
        "--go",
        default="routes/routes.go",
        help="Path to routes.go  (default: routes/routes.go)",
    )
    parser.add_argument(
        "--swagger",
        default="docs/swagger.yaml",
        help="Path to swagger.yaml  (default: docs/swagger.yaml)",
    )
    args = parser.parse_args()

    go_path = Path(args.go)
    sw_path = Path(args.swagger)

    missing = [p for p in (go_path, sw_path) if not p.exists()]
    if missing:
        for p in missing:
            print(f"ERROR: file not found: {p}", file=sys.stderr)
        sys.exit(2)

    sys.exit(compare_and_report(go_path, sw_path))