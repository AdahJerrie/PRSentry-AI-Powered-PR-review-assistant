# PRSentry — API Contract

This is the locked contract for data crossing service boundaries. All three services build against this without needing each other's code to exist yet. If you need to change a field here, flag it to the team before changing your own service's code — a silent change on one side breaks the other two.

---

## 1. Go → Python: Diff Analysis

**Endpoint:** `POST /review`

**Request (Go → Python)**

```json
{
  "pr_id": 123,
  "repo": "owner/name",
  "files": [
    {
      "path": "main.go",
      "diff": "@@ -10,6 +10,7 @@ ... unified diff ..."
    }
  ]
}
```

| Field         | Type            | Notes                                                        |
|---------------|-----------------|---------------------------------------------------------------|
| `pr_id`       | int             | GitHub PR number, not a global ID                              |
| `repo`        | string          | `"owner/name"` format                                          |
| `files`       | array of objects| One entry per changed file                                     |
| `files[].path`| string          | File path relative to repo root                                 |
| `files[].diff`| string          | Unified diff text for that file only                             |

**Response (Python → Go)**

```json
{
  "pr_id": 123,
  "summary": "Refactors error handling in the auth package...",
  "risk_score": 3,
  "findings": [
    {
      "file_path": "main.go",
      "line_number": 42,
      "severity": "medium",
      "category": "bug",
      "message": "Error from db.Query is ignored — could mask connection failures."
    }
  ]
}
```

| Field                    | Type   | Notes                                                                 |
|--------------------------|--------|------------------------------------------------------------------------|
| `pr_id`                  | int    | Echoed back so Go can match the response to the right PR                 |
| `summary`                | string | Human-readable overview of the PR                                        |
| `risk_score`             | number | Numeric scale (e.g. 1–5). Not an enum string — lets Go threshold on it     |
| `findings`               | array  | Can be an empty array — a clean PR is a completed review with zero findings, not an error |
| `findings[].file_path`   | string | Matches a `files[].path` from the request                                |
| `findings[].line_number` | int    | Line in the file the finding applies to                                   |
| `findings[].severity`    | string | Fixed set: `"low"`, `"medium"`, `"high"`                                  |
| `findings[].category`    | string | Fixed set: `"bug"`, `"security"`, `"style"`, `"maintainability"`          |
| `findings[].message`     | string | Human-readable explanation of the finding                                 |

**v1 guardrail — large PRs:** If a PR exceeds a file-count threshold (Go-side constant, currently uncommitted to a specific number), Go skips calling Python entirely and either truncates the file list or returns a "PR too large to auto-review" message. This is intentionally simple for v1; smarter chunking/prioritization is deferred to v2.

---

## 2. GitHub → Go: Webhook Events

Go listens for `pull_request` events at `POST /webhook`, verified via HMAC-SHA256 signature (`X-Hub-Signature-256` header) against a shared webhook secret.

**Actions Go acts on:** `opened`, `synchronize`, `reopened`
**Actions Go acknowledges but ignores:** everything else (e.g. `closed`) — still returns `200 OK` so GitHub doesn't retry, but no review is triggered.

---

## 3. Go → React: Dashboard API

| Endpoint                              | Method | Purpose                                      |
|----------------------------------------|--------|-----------------------------------------------|
| `/api/prs?repo=&status=&limit=`        | GET    | List of PRs with latest review status          |
| `/api/prs/:id`                         | GET    | PR detail + its review                          |
| `/api/reviews/:id`                     | GET    | Full review with all findings                    |
| `/api/stats`                           | GET    | Aggregate counts for the trends page              |

Response shapes for these are **not yet locked** — finalize once Person A's `internal/api` package and Person C's dashboard data needs are both clearer (targeted for Week 2 per the build plan).

---

## Change Log

| Date       | Change                          | Reason                                      |
|------------|----------------------------------|----------------------------------------------|
| (fill in)  | Initial contract locked          | Week 1 kickoff                                |
