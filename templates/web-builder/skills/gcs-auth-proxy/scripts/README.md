# gcs-auth-proxy

A minimal Cloud Run service that proxies authenticated requests to a private Google Cloud Storage bucket. It allows you to serve static assets from GCS through Google IAP (Identity-Aware Proxy) without making the bucket public.

## How it works

1. Cloud Run enforces Google authentication via IAP (`roles/run.invoker`)
2. The service account running the Cloud Run service has `roles/storage.objectViewer` on the bucket
3. The proxy strips a configurable URL prefix, fetches the corresponding object from GCS, and streams it back with correct content-type and cache headers

This means the GCS bucket stays private — all access is through the authenticated Cloud Run URL.

## Configuration

All configuration is via environment variables:

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `GCS_BUCKET` | **Yes** | — | GCS bucket name to serve from |
| `STRIP_PREFIX` | No | `/<GCS_BUCKET>` | URL prefix to strip before mapping to a GCS object name |
| `ROOT_REDIRECT` | No | `/index.html` | Where to redirect bare `/` requests |
| `PORT` | No | `8080` | Port to listen on (Cloud Run sets this automatically) |

### STRIP_PREFIX examples

| Request path | STRIP_PREFIX | GCS object |
|---|---|---|
| `/my-bucket/assets/logo.png` | `/my-bucket` | `assets/logo.png` |
| `/assets/logo.png` | *(empty)* | `assets/logo.png` |
| `/site/en/index.html` | `/site` | `en/index.html` |

The default (`/<GCS_BUCKET>`) works naturally when your frontend sets its base URL to `/<bucket-name>/...` — requests already include the bucket name in the path.

## Deploy to Cloud Run

```bash
# Build and push
gcloud builds submit --tag gcr.io/YOUR_PROJECT/gcs-auth-proxy

# Deploy
gcloud run deploy gcs-auth-proxy \
  --image gcr.io/YOUR_PROJECT/gcs-auth-proxy \
  --region us-central1 \
  --no-allow-unauthenticated \
  --set-env-vars GCS_BUCKET=your-bucket-name \
  --set-env-vars STRIP_PREFIX=/your-bucket-name \
  --set-env-vars ROOT_REDIRECT=/your-bucket-name/index.html \
  --service-account your-sa@your-project.iam.gserviceaccount.com
```

Grant the service account access to the bucket:

```bash
gcloud storage buckets add-iam-policy-binding gs://your-bucket-name \
  --member="serviceAccount:your-sa@your-project.iam.gserviceaccount.com" \
  --role="roles/storage.objectViewer"
```

Restrict Cloud Run invocation to your org (optional but recommended):

```bash
gcloud run services add-iam-policy-binding gcs-auth-proxy \
  --region us-central1 \
  --member="domain:your-domain.com" \
  --role="roles/run.invoker"
```

## Local development

```bash
# Authenticate with your Google account
gcloud auth application-default login

# Run locally
GCS_BUCKET=your-bucket-name \
STRIP_PREFIX=/your-bucket-name \
go run main.go
```

## Build

```bash
docker build -t gcs-auth-proxy .
docker run -p 8080:8080 \
  -e GCS_BUCKET=your-bucket-name \
  -e GOOGLE_APPLICATION_CREDENTIALS=/path/to/key.json \
  gcs-auth-proxy
```

## Features

- Streams objects directly from GCS — no local buffering
- Preserves `Content-Type`, `Cache-Control`, and `Content-Encoding` headers from GCS object metadata
- MIME type inference from file extension as fallback when GCS metadata is missing
- HTTP/2 (h2c) support for Cloud Run
- Zero-dependency runtime image (distroless)
