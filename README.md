# 🏛️ Congress CLI

A terminal-based tool to explore U.S. congressional activity — including bills, votes, and roll call summaries — from both Congress.gov and the Clerk of the House.

Built in Go for clarity, speed, and transparency.

---

## 📦 Features

- 🗳️ List recent bills in a given Congress
- 🔍 Search bills by ID and view their vote history
- 📜 Show detailed roll call results from the Clerk's XML feed
- 🧑‍🤝‍🧑 See how each party voted
- ✅ Gracefully handles gaps in Congress.gov API

---

## 🚀 Getting Started

### Prerequisites

- Go 1.20 or later
- A [Congress.gov API Key](https://api.congress.gov/)

### Setup

```bash
git clone https://github.com/your-username/congresscli.git
cd congresscli

# Add your API key to .env
echo "CONGRESS_API_KEY=your_api_key_here" > .env

# Install Go dependencies
go mod tidy

🛠️ Usage


List recent bills
```
go run . list --congress=119
```
        


Show vote results for a bill
```
go run . votes --congress=119 --id hr2670
```

If Congress.gov doesn’t return vote data, the CLI will prompt you to fall back to Clerk.gov roll calls.

Directly query a known roll call
```
go run . clerkvote --year=2023 --roll=328
```

💡 Example Output
```
🧮 Vote Totals by Party:
• Republican: 🟢 215 | 🔴 4 | ⚪ 0 | ❌ 3
• Democratic: 🟢 4   | 🔴 206 | ⚪ 0 | ❌ 2
• Independent: 🟢 0 | 🔴 0 | ⚪ 0 | ❌ 0
```

🔧 Project Structure
```
congresscli/
├── main.go          # CLI entrypoint
├── bill.go          # Handles Congress.gov bill/vote logic
├── votes.go         # Handles Clerk XML roll call parsing
├── .env             # Stores your Congress.gov API key
└── go.mod           # Module definition
```

📓 Planned Features
- Filter bills by chamber, sponsor, or status

- List all votes per bill (not just the most recent)

- Export votes to CSV/JSON

- Auto-open bill or roll call in browser

🤝 Contributing
Pull requests and feature ideas welcome! Open an issue or fork and submit a PR.


