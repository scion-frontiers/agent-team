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
human a link. Commands below use the bare-uppercase placeholders `BUCKET_NAME`, `PROJECT_ID`,
`PREFIX`, `OBJECT_PATH` and `PROJECT_FOLDER`; substitute your own values before running anything.
In the Scion hub environment the defaults are bucket `ddt-scion-hub-exchange` in project
`deploy-demo-test`.

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
  Objects are served at `https://storage.googleapis.com/BUCKET_NAME/OBJECT_PATH` — **no per-object ACL
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

**Always pass both `-h "Content-Type:…"` and `-h "Cache-Control:…"` explicitly, on every upload,
including HTML.** Both header values contain spaces and commas, so each `-h` argument must stay
quoted exactly as shown.

```bash
# HTML page -> becomes https://storage.googleapis.com/BUCKET_NAME/PREFIX/index.html
gsutil -h "Content-Type:text/html" \
       -h "Cache-Control:no-cache, no-store, must-revalidate" \
       cp mypage.html gs://BUCKET_NAME/PREFIX/index.html

# Source data alongside it
gsutil -h "Content-Type:text/yaml" \
       -h "Cache-Control:no-cache, no-store, must-revalidate" \
       cp data.yaml gs://BUCKET_NAME/PREFIX/data.yaml
```

Notes:
- **Set `Content-Type` yourself. Do not rely on it being inferred.** When the type is not set,
  the object can be served as `application/octet-stream`, and a browser then **downloads the
  file instead of rendering it** — the page appears broken to the reviewer while the upload
  itself reported success. This is the most common way a published artifact fails. Common
  values: `text/html`, `text/css`, `application/javascript`, `text/yaml`, `image/png`,
  `image/svg+xml`.
- **`Cache-Control:no-cache, no-store, must-revalidate`** is important — without it, re-uploads
  (re-renders) may not show up for the reviewer due to CDN/browser caching. This is the case
  every time a human is iterating on an artifact you keep republishing to the same URL.
- Naming the entry page `index.html` gives a clean link (`…/PREFIX/index.html`).

---

## 3. Verify it's live and public

**Verify against the build you just uploaded, not against the URL.** A `200`, a `text/html`, or
a grep for a string the page has always contained will all pass **before** your upload as
happily as after it — they confirm that *something* is being served, which was already true.

> **A verification that would have passed before the upload is not a verification.**

So give each build a marker that could not have been there a minute ago, and check for **that**:

```bash
BUILD_ID="build-$(date -u +%Y%m%dT%H%M%SZ)"

# 1. Stamp the marker into the artifact BEFORE uploading — from your renderer, or appended:
printf '<!-- %s -->\n' "$BUILD_ID" >> page.html

# 2. Upload
gsutil -h "Content-Type:text/html" \
       -h "Cache-Control:no-cache, no-store, must-revalidate" \
       cp page.html gs://BUCKET_NAME/PREFIX/index.html

# 3. Fetch what is actually being served and look for THIS build's marker
body=$(curl -fsS "https://storage.googleapis.com/BUCKET_NAME/PREFIX/index.html"); fetch_rc=$?
printf '%s' "$body" | grep -q "$BUILD_ID"; marker_rc=$?
echo "fetch_rc=$fetch_rc marker_rc=$marker_rc"   # both 0 -> this build is the one being served
```

`fetch_rc` and `marker_rc` answer two different questions and both matter: the first is whether
the object is reachable and public at all, the second is whether the bytes coming back are the
ones you just wrote. **Capture each separately** — in a pipeline, `$?` is only the last command's
status.

Headers are worth checking too, but check them as a **separate** claim; they are metadata and say
nothing about which build is live:

```bash
gsutil stat gs://BUCKET_NAME/PREFIX/index.html | grep -iE "Content-Type|Cache"
curl -sSI "https://storage.googleapis.com/BUCKET_NAME/PREFIX/index.html" \
  | grep -iE "^(content-type|cache-control):"
```

If you have no way to stamp a marker, the fallback is to grep for a string **generated by this
run's content** — a new row, a new figure, a timestamp the renderer wrote. Never a title, a
heading, or a nav label; those survive every rebuild and make the check decorative.

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
cd PROJECT_FOLDER
BUILD_ID="build-$(date -u +%Y%m%dT%H%M%SZ)"

python3 render.py data.yaml page.html                 # your renderer
printf '<!-- %s -->\n' "$BUILD_ID" >> page.html       # marker for the verification below

gsutil -h "Content-Type:text/html" \
       -h "Cache-Control:no-cache, no-store, must-revalidate" \
       cp page.html gs://BUCKET_NAME/PREFIX/index.html
gsutil -h "Content-Type:text/yaml" \
       -h "Cache-Control:no-cache, no-store, must-revalidate" \
       cp data.yaml gs://BUCKET_NAME/PREFIX/data.yaml

body=$(curl -fsS "https://storage.googleapis.com/BUCKET_NAME/PREFIX/index.html"); fetch_rc=$?
printf '%s' "$body" | grep -q "$BUILD_ID"; marker_rc=$?
echo "fetch_rc=$fetch_rc marker_rc=$marker_rc"        # both 0 before you share the link
# -> share https://storage.googleapis.com/BUCKET_NAME/PREFIX/index.html
```

## Gotchas recap

- Public bucket — no secrets, ever.
- Use a per-effort prefix; check existing prefixes first.
- `Cache-Control:no-cache, no-store, must-revalidate` so re-renders are visible.
- Set `Content-Type` explicitly on **every** upload, HTML included — unset means
  `application/octet-stream`, which downloads instead of rendering.
- Verify by grepping the live page for a marker unique to **this** build. A check that would have
  passed before the upload is not a check.
- Self-contained HTML (no CDN); robust YAML loading via ruamel fallback.
