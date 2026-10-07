# Country Game Board Generator

A Go web app that generates one-page, legal-size printable PDF game boards.

## Current controls

- Country: 195 sovereign countries
- Difficulty: Beginner, Easy, Medium, Hard, Insane
- Territories: 6 through 45
- No combination tracking or locking

## PDF

The PDF currently includes:

- Country name
- Country flag (downloaded from FlagCDN when available)
- Country area in km²
- Difficulty
- Territory count
- Territory grid
- Territory number
- Territory payout
- Required roll
- Number of dice
- Country conquest reward
- Legend

## Important rule note

The supplied country-reward table defines rewards only at 5, 10, 15, ... 45 territories. Because the app now permits 6-45 territories, the first version uses the **next 5-territory bracket** for the country reward. For example, 18 territories uses the 20-territory reward. This logic is isolated in `countryReward()` so it can be changed after testing.

## Territory names

The first build uses `Territory 1`, `Territory 2`, etc. This is intentional for the first deployment/test pass. Country-specific province/city/territory names can be added after the PDF layout is confirmed.

## Run locally

```powershell
go mod tidy
go run .
```

Then open:

http://localhost:8080

## Render

Build command:

```text
go build -tags netgo -ldflags '-s -w' -o app
```

Start command:

```text
./app
```

Render automatically supplies the `PORT` environment variable, and the application binds to `0.0.0.0`.
