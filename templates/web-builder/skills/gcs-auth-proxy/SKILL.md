# Skill: GCS Auth Proxy

This skill describes the architecture and deployment pattern for serving private Google Cloud Storage (GCS) bucket content behind Identity-Aware Proxy (IAP) authentication using a Cloud Run reverse proxy.

## Architecture

This pattern allows you to lock down a GCS bucket (removing all public access) while still serving its content to authenticated users.

1.  **Storage**: A GCS bucket contains the static site assets. Public access is disabled.
2.  **Proxy**: A lightweight HTTP server (typically written in Go) runs on Cloud Run.
3.  **Authentication**: IAP is enabled on the Cloud Run service, requiring users to authenticate before they can reach the proxy.
4.  **Flow**: User -> Cloud Run (IAP Auth) -> Proxy Server -> GCS (Internal Request).

The proxy server uses a dedicated Service Account with `roles/storage.objectViewer` permissions on the bucket to fetch and serve the requested content.

## Deployment Pattern

### 1. Service Account Setup
Create a service account specifically for the proxy and grant it access to the bucket.

```bash
gcloud iam service-accounts create SERVICE_ACCOUNT_NAME \
    --display-name="GCS Auth Proxy Service Account"

gsutil iam ch serviceAccount:SERVICE_ACCOUNT_NAME@PROJECT_ID.iam.gserviceaccount.com:roles/storage.objectViewer \
    gs://BUCKET_NAME
```

### 2. Cloud Run Deployment
Deploy the proxy to Cloud Run. Ensure unauthenticated access is disabled to allow IAP to function.

```bash
gcloud run deploy SERVICE_NAME \
    --source . \
    --region REGION \
    --no-allow-unauthenticated \
    --set-env-vars GCS_BUCKET=BUCKET_NAME \
    --service-account SERVICE_ACCOUNT_NAME@PROJECT_ID.iam.gserviceaccount.com
```

*Note: The "Setting IAM policy failed" warning about `allUsers` is expected and can be ignored when using `--no-allow-unauthenticated`.*

### 3. IAP Enablement and Access Control
Enable IAP on the service and grant access to authorized users.

```bash
gcloud iap web add-iam-policy-binding \
    --resource-type=cloud-run \
    --service=SERVICE_NAME \
    --region=REGION \
    --member="user:USER_EMAIL" \
    --role="roles/iap.httpsResourceAccessor"
```

## Implementation Lessons

When implementing the proxy server, keep these battle-tested lessons in mind:

- **HTTP/2 Support**: Cloud Run often uses HTTP/2 (H2C) to communicate with the container. Ensure your server can handle H2C. In Go, you may need to wrap your handler with `h2c.NewHandler`.
- **MIME Type Fallback**: GCS metadata sometimes lacks accurate `Content-Type` information (often defaulting to `text/plain` or `application/octet-stream`). Implement a fallback mechanism that uses file extensions to determine the correct MIME type.
- **Relative Pathing**: All resource references (CSS, JS, Images, Video) must use relative proxy paths. Do not link to absolute `storage.googleapis.com` URLs, as those will return 403 errors once the bucket is secured.
- **Build Asset Hashing**: Ensure that all related assets (HTML, CSS, JS) are uploaded as a complete set. Incomplete uploads or mismatched hashes will lead to 404 errors that browsers may misleadingly report as MIME type mismatches.
- **JWT Audience**: The JWT audience for IAP on Cloud Run follows this pattern: `/projects/PROJECT_NUMBER/locations/REGION/services/SERVICE_NAME`. This is different from load-balancer-based IAP configurations.

## Reference Implementation

A complete, battle-tested Go implementation is included in the `scripts/` directory alongside this skill:

- **`scripts/main.go`** — The proxy server. Handles GCS object fetching, MIME type fallback, H2C support, configurable path stripping, and root redirect.
- **`scripts/Dockerfile`** — Multi-stage build producing a minimal distroless image.
- **`scripts/go.mod` / `scripts/go.sum`** — Go module dependencies.
- **`scripts/README.md`** — Configuration reference and deployment commands.

### Configuration (Environment Variables)

| Variable | Required | Default | Description |
|---|---|---|---|
| `GCS_BUCKET` | **Yes** | — | GCS bucket name to serve from |
| `STRIP_PREFIX` | No | `/<GCS_BUCKET>` | URL prefix to strip before mapping to GCS object name |
| `ROOT_REDIRECT` | No | `/index.html` | Where to redirect bare `/` requests |
| `PORT` | No | `8080` | Listen port (Cloud Run sets this automatically) |

To deploy, copy the `scripts/` directory contents into your project, then use `gcloud builds submit` and `gcloud run deploy` as described in the Deployment Pattern section above.
