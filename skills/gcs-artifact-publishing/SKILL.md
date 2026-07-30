---
name: gcs-artifact-publishing
description: >-
  Publish human-reviewable artifacts (HTML dashboards, reports, rendered docs) to the shared
  Scion GCS artifact bucket and share a public link. Covers the bucket/auth setup, public-read
  serving model, per-effort prefix isolation, gsutil upload with correct Content-Type and
  Cache-Control, verifying public access, and building self-contained HTML renderers (including
  the no-PyYAML fallback to the gcloud-bundled ruamel). Use whenever an agent needs to give a
  human a clickable link to review generated output instead of pasting it into chat.
---

# GCS Artifact Publishing — Scion

A practical reference for publishing reviewable artifacts to the shared bucket and handing a
human a link. Commands below use `BUCKET_NAME` and `PROJECT_ID`; substitute the bucket and
project your environment publishes to. In the Scion hub environment the defaults are bucket
`ddt-scion-hub-exchange` in project `deploy-demo-test`.

## When to use

- You generated something a human should review visually (an HTML dashboard, a report, a rendered
  roadmap/plan) and chat/markdown is a poor medium.
- You want a stable, shareable URL rather than a file path on a volume only agents can see.

## Prerequisites

- `gcloud` / `gsutil` authenticated (in agent sandboxes this is typically a service account, e.g.
  `SERVICE_ACCOUNT_NAME@PROJECT_ID.iam.gserviceaccount.com`, with `storage.admin` on the bucket).
- `python3` if you are rendering HTML from data.

Verify access:
```bash
gcloud auth list
gsutil ls gs://BUCKET_NAME/
```

---

## 1. Serving model (read this first)

- The artifact bucket must be **public-read** via bucket IAM (`allUsers` has `roles/storage.objectViewer`).
  Objects are served at `https://storage.googleapis.com/<bucket>/<path>` — **no per-object ACL
  needed**. Confirm with:
  ```bash
  gcloud storage buckets get-iam-policy gs://BUCKET_NAME --format=json | grep -A2 allUsers
  ```
- **Anything you upload to a public-read bucket is world-readable.** Never publish secrets, tokens, internal
  credentials, or sensitive customer data.
- **Isolate each effort under its own prefix** (e.g. `roadmap/`, `assetview/`) so concurrent
  efforts don't collide. List existing prefixes before picking one:
  ```bash
  gsutil ls gs://BUCKET_NAME/
  ```

---

## 2. Uploading

```bash
# HTML page -> becomes https://storage.googleapis.com/BUCKET_NAME/<prefix>/index.html
gsutil -h "Cache-Control:no-cache, max-age=0" \
  cp mypage.html gs://BUCKET_NAME/<prefix>/index.html

# Source data alongside it (set Content-Type for non-HTML so browsers render/serve correctly)
gsutil -h "Cache-Control:no-cache, max-age=0" -h "Content-Type:text/yaml" \
  cp data.yaml gs://BUCKET_NAME/<prefix>/data.yaml
```

Notes:
- **`Cache-Control:no-cache`** is important — without it, re-uploads (re-renders) may not show up
  for the reviewer due to CDN/browser caching.
- `gsutil` infers `Content-Type` from the extension for common types (`.html` → `text/html`).
  Set it explicitly for `.yaml`/`.md`/etc. if you want a specific type.
- Naming the entry page `index.html` gives a clean link (`…/<prefix>/index.html`).

---

## 3. Verify it's live and public

```bash
gsutil stat gs://BUCKET_NAME/<prefix>/index.html | grep -iE "Content-Type|Cache"
curl -s -o /dev/null -w "%{http_code} %{content_type}\n" \
  "https://storage.googleapis.com/BUCKET_NAME/<prefix>/index.html"   # expect: 200 text/html
```

Then share the URL with the human (on their channel — see `scion-messaging`).

---

## 4. Building self-contained HTML renderers

For review pages, prefer **self-contained HTML**: inline CSS + vanilla JS, with the data embedded
as a JSON literal. No external/CDN dependencies means the page renders straight from the bucket
with nothing else to fetch. Pattern: read source data → `json.dumps` it into an HTML template
string → write file → upload. A little client-side JS (group-by toggle, filters) makes large lists
reviewable without server support.

### YAML loading gotcha (agent sandboxes)

Agent sandboxes often have **no system PyYAML and no `pip`**. The gcloud SDK vendors a PyYAML that
won't import standalone (py2-style imports), **but its bundled `ruamel.yaml` works**. Make renderers
robust with this fallback:

```python
def load_yaml(path):
    try:
        import yaml                      # normal case
        with open(path) as f: return yaml.safe_load(f)
    except ImportError:
        import os, sys
        for c in ["/usr/lib/google-cloud-sdk/lib/third_party",
                  os.path.expanduser("~/google-cloud-sdk/lib/third_party")]:
            if os.path.isdir(os.path.join(c, "ruamel")):
                sys.path.insert(0, c); break
        from ruamel import yaml as ry
        y = ry.YAML(typ="safe")
        with open(path) as f: return y.load(f)
```

(Find the gcloud-bundled interpreter with `gcloud info --format="value(basic.python_location)"`.)
No Ruby is available and Node lacks `js-yaml`; for quick validation without a parser, fall back to
structural string checks in Python.

---

## 5. End-to-end example

```bash
cd <project-folder>
python3 render.py data.yaml page.html          # your renderer
gsutil -h "Cache-Control:no-cache, max-age=0" cp page.html gs://BUCKET_NAME/<prefix>/index.html
gsutil -h "Cache-Control:no-cache, max-age=0" -h "Content-Type:text/yaml" cp data.yaml gs://BUCKET_NAME/<prefix>/data.yaml
curl -s -o /dev/null -w "%{http_code}\n" "https://storage.googleapis.com/BUCKET_NAME/<prefix>/index.html"
# -> share https://storage.googleapis.com/BUCKET_NAME/<prefix>/index.html
```

## Gotchas recap

- Public bucket — no secrets, ever.
- Use a per-effort prefix; check existing prefixes first.
- `Cache-Control:no-cache` so re-renders are visible.
- Set `Content-Type` for non-HTML files.
- Self-contained HTML (no CDN); robust YAML loading via ruamel fallback.
