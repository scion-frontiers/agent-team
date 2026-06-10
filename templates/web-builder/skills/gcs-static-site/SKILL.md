# Skill: GCS Static Site Publishing

This skill covers the technical workflow for publishing and maintaining static websites on Google Cloud Storage (GCS).

## Publishing to GCS

When publishing files to a GCS bucket, you must ensure they are served with the correct metadata so browsers can render them properly.

### The Publishing Command
Always use `gsutil` with explicit headers for `Content-Type` and `Cache-Control`.

```bash
gsutil -h "Content-Type:text/html" \
       -h "Cache-Control:no-cache, no-store, must-revalidate" \
       cp LOCAL_FILE_PATH gs://BUCKET_NAME/TARGET_PATH
```

### Key Technical Details
- **Content-Type**: Without this header, GCS may default to `application/octet-stream`, causing browsers to download the file instead of rendering it. Common types:
  - `text/html`
  - `text/css`
  - `application/javascript`
  - `image/png`
  - `image/svg+xml`
- **Cache-Control**: Use `no-cache, no-store, must-revalidate` for sites that update frequently. This prevents users from seeing stale content and avoids the "it works in incognito but not normally" debugging pain.
- **Verify with Curl**: Never assume a publish was successful. Verify by checking the live URL:
  ```bash
  curl -s -I "https://storage.googleapis.com/BUCKET_NAME/TARGET_PATH" | grep "Content-Type"
  curl -s "https://storage.googleapis.com/BUCKET_NAME/TARGET_PATH" | grep "UNIQUE_STRING_FROM_YOUR_UPDATE"
  ```
- **No Scion Storage Command**: All storage operations must use `gsutil`.

## Managing Bucket Contents

Use `gsutil` to audit and manage the artifacts in your bucket.

- **List Contents**: `gsutil ls gs://BUCKET_NAME/`
- **Recursive Listing**: `gsutil ls -r gs://BUCKET_NAME/`
- **Search for Artifacts**: `gsutil ls -r gs://BUCKET_NAME/ | grep -i PATTERN`

## Batch Processing Updates

To optimize the workflow and reduce GCS writes, batch multiple updates when they arrive in quick succession.

- **Optimal Batch size**: Process 2-3 notifications or updates into a single local edit and upload cycle.
- **Avoid Over-batching**: Do not wait too long; new updates arriving mid-edit can cause merge conflicts or complexity in your local state.
- **Immediate Processing**: If updates are infrequent, process and publish each one immediately to keep the site current.

## Common Pitfalls

- **Forgetting Headers**: Always double-check that you included the `-h` flags in your `gsutil cp` command.
- **Browser Caching**: If a user reports that they don't see updates, suggest a hard refresh or checking in an incognito window, and verify your `Cache-Control` headers.
- **Count Drift**: When updating counters (e.g., "Total Modules: 86"), ensure you've read the latest version of the file to avoid off-by-one errors or overwriting other agents' updates.
- **Partial Uploads**: Always upload all related build assets (HTML, CSS, JS) together. Partial uploads can lead to broken references or MIME type errors reported by browsers.

## Site Design Lessons

- **Curate, Don't Catalog**: Your goal is to tell a story of project progress. Avoid listing every single minor artifact. Focus on milestones and significant outcomes.
- **Effective Page Structure**:
  1. **Hero Section**: High-level stats and bold statements.
  2. **The Narrative**: A chronological or thematic story of the project's evolution.
  3. **Live Links**: Direct access to running applications or major deliverables.
  4. **Collapsible Details**: Use the HTML `<details>` and `<summary>` tags to keep the page clean while allowing stakeholders to drill down into specifics.
- **Native Interactivity**: Prefer native HTML features like `<details>` for interactivity to keep the site lightweight and avoid unnecessary JavaScript.
