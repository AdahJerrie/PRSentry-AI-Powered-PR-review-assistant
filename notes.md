Dimension	        Complexity	    Risk Level	            Mitigation Strategy

Development	        Moderate	        Low	               Mocking endpoints early allows continuous integration.

GitHub Integration	High	            Medium	           App permissions and Webhook handling have a steep initial learning curve.

LLM Reliability	    Moderate	    Medium	           Rely strictly on Structured Outputs (JSON mode) to prevent UI parsing crashes.

Deployment	        Low	                Low	            Dockerization hides OS differences; platforms like Railway make hosting simple.

# Some of the bottle necks I was looking into
1. The "Large Diff" Problem (Python Track)
GitHub diffs can easily exceed LLM context limits or trigger massive API token bills if a PR alters lock files (e.g., package-lock.json, go.sum) or auto-generated assets.
    *Action: Implement an early file-filtering layer in Go or Python. Immediately ignore binary files, lockfiles, and minified assets before passing the payload to the LLM.*

2. GitHub Webhook Timeouts (Go Track)GitHub expects your webhook endpoint to acknowledge a payload with a 200 OK response within 10 seconds. Calling the Python service, waiting for the LLM to complete its analysis, and returning the data synchronously will frequently cross this 10-second threshold for larger PRs.
    *Action: If a synchronous architecture causes GitHub timeout errors during Week 2, immediately pivot Go's webhook handler to save the incoming event to the database with a pending status, return a 202 Accepted to GitHub, and launch a Go goroutine (go internal.AnalyzePR(...)) to process the LLM call asynchronously in the background.*

3. Inline Comment Mapping
Posting an overall summary comment on a PR is trivial. However, posting inline findings tied to exact lines (findings.line_number) requires understanding GitHub's specific review comment API, which maps to the diff's relative line index (the "position" field in the diff hunk), not necessarily the absolute file line number.
    *Action: Dedicate extra testing time in Week 2 for the mapping logic, or fall back to posting a single beautifully formatted Markdown table of findings as a main PR comment if inline mapping becomes a blocker.*