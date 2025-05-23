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
	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)
	searchID := searchCmd.String("id", "", "Bill ID (e.g., hr2250)")
	searchKeyword := searchCmd.String("keyword", "", "Keyword to search for in bill titles")

	clerkCmd := flag.NewFlagSet("clerkvote", flag.ExitOnError)
	clerkYear := clerkCmd.String("year", "", "Year of the vote (e.g. 2023)")
	clerkRoll := clerkCmd.String("roll", "", "Roll call number (e.g. 328)")
	
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
			fmt.Println("Provide a bill ID with --id")
			return
		}
		congress.FetchVotesByBillID(apiKey, *votesID)

	case "list":
		_ = listCmd.Parse(os.Args[2:])
		congress.FetchRecentBills(apiKey)

	case "search":
		_ = searchCmd.Parse(os.Args[2:])
		if *searchID != "" {
			congress.FetchBillByID(apiKey, *searchID)
		} else if *searchKeyword != "" {
			congress.SearchBillsByKeyword(apiKey, *searchKeyword)
		} else {
			fmt.Println("Provide either --id or --keyword")
		}

	default:
		fmt.Println("expected 'list' or 'search' subcommand")
		os.Exit(1)
	}
}
