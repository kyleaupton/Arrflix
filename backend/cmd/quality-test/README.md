# quality-test

A CLI tool that compares Arrflix's quality parser against a live Sonarr instance to validate parity.

## Purpose

This tool queries Sonarr's `/api/v3/parse` endpoint with test release titles and compares the quality detection results against Arrflix's internal `quality.ParseQuality()` function. It helps ensure our Go port of Sonarr's quality parsing logic returns identical results.

## Usage

```bash
SONARR_URL=http://localhost:8989 SONARR_API_KEY=your-api-key go run ./cmd/quality-test
```

## Environment Variables

| Variable         | Description                                                      |
| ---------------- | ---------------------------------------------------------------- |
| `SONARR_URL`     | Full URL to your Sonarr instance (e.g., `http://localhost:8989`) |
| `SONARR_API_KEY` | Your Sonarr API key (found in Settings → General)                |

## Output

The tool displays:

- Each test title with match/mismatch status
- Summary statistics (matches, mismatches, parity percentage)
- Detailed breakdown of any mismatches showing both Sonarr's and Arrflix's results

## Related

- **Unit tests**: `backend/internal/quality/quality_test.go` contains the same test cases with hardcoded expected values for fast CI testing without needing a Sonarr instance.
- **Quality parser**: `backend/internal/quality/quality.go` contains the actual parsing logic.
