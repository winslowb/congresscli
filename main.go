package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"congresscli/congress"
	"github.com/joho/godotenv"
)

func main() {
	
	votesCmd := flag.NewFlagSet("votes", flag.ExitOnError)
	votesID := votesCmd.String("id", "", "Bill ID (e.g., hr1818)")

	_ = godotenv.Load()
	apiKey := os.Getenv("CONGRESS_API_KEY")
	if apiKey == "" {
		log.Fatal("CONGRESS_API_KEY not set")
	}

	// Subcommands
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	votesCongress := votesCmd.String("congress", "", "Congress number (e.g., 119)")

	clerkCmd := flag.NewFlagSet("clerkvote", flag.ExitOnError)
	clerkYear := clerkCmd.String("year", "", "Year of the vote (e.g. 2023)")
	clerkRoll := clerkCmd.String("roll", "", "Roll call number (e.g. 328)")


	listCongress := listCmd.String("congress", "119", "Congress number (e.g., 118, 117)")
	
	if len(os.Args) < 2 {
		fmt.Println("expected 'list' or 'search' subcommand")
		os.Exit(1)
	}

var session string
flag.StringVar(&session, "session", "1", "Session number (1 or 2)")

	switch os.Args[1] {


	case "clerkvote":
		_ = clerkCmd.Parse(os.Args[2:])
		if *clerkYear == "" || *clerkRoll == "" {
			fmt.Println("Provide both --year and --roll")
			return
		}
		congress.FetchClerkXMLRollCall(*clerkYear, *clerkRoll)

	case "votes":
		_ = votesCmd.Parse(os.Args[2:])
		if *votesID == "" {
			fmt.Println("Provide both --id and --congress")
			return
		}
		congress.FetchVotesByBillID(apiKey, *votesID, *votesCongress)

	case "list":
		_ = listCmd.Parse(os.Args[2:])
		err := congress.FetchRecentBills(apiKey, *listCongress)
		if err != nil {
			log.Fatalf("Failed to fetch bills: %v", err)
		}

	default:
		fmt.Println("expected 'list' or 'search' subcommand")
		os.Exit(1)
	}
}
