package congress 

import (
	"encoding/json"
//	"os"
	"fmt"
	"net/http"
	"strings"
	"regexp"
)

// Fetch specific bill by ID (e.g., hr2250)

func FetchBillByID(apiKey, billID string) {
	parts := strings.Split(billID, "r")
	if len(parts) != 2 {
		fmt.Println("Invalid bill ID format. Example: hr2250")
		return
	}
	billType := "hr"
	billNumber := parts[1]

	url := fmt.Sprintf("https://api.congress.gov/v3/bill/119/%s/%s?format=json", billType, billNumber)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Api-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching bill:", err)
		return
	}
	defer resp.Body.Close()

	var result struct {
		Bill struct {
			Title              string `json:"title"`
			Number             string `json:"number"`
			OriginChamber      string `json:"originChamber"`
			IntroducedDate     string `json:"introducedDate"`
			Type               string `json:"type"`
			PolicyArea         struct{ Name string } `json:"policyArea"`
			LatestAction       struct {
				ActionDate string `json:"actionDate"`
				Text       string `json:"text"`
			} `json:"latestAction"`
			Sponsors []struct {
				FullName string `json:"fullName"`
				Party    string `json:"party"`
				State    string `json:"state"`
			} `json:"sponsors"`
			TextVersions struct {
				Count int    `json:"count"`
				URL   string `json:"url"`
			} `json:"textVersions"`
			Titles struct {
				URL string `json:"url"`
			} `json:"titles"`
		} `json:"bill"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Decode error:", err)
		return
	}

	b := result.Bill
	fmt.Printf("\n📄 %s %s: %s\n", b.Type, b.Number, b.Title)
	fmt.Printf("🏛️  Chamber: %s\n", b.OriginChamber)
	fmt.Printf("📅 Introduced: %s\n", b.IntroducedDate)
	if len(b.Sponsors) > 0 {
		s := b.Sponsors[0]
		fmt.Printf("🧑 Sponsor: %s (%s-%s)\n", s.FullName, s.Party, s.State)
	}
	fmt.Printf("🗂️  Policy Area: %s\n", b.PolicyArea.Name)
	fmt.Printf("📌 Status: %s on %s\n", b.LatestAction.Text, b.LatestAction.ActionDate)

	// External links
	fmt.Println("\n🔗 Resources:")
	fmt.Printf("• Summary: %s\n", b.Titles.URL)
	fmt.Printf("• Full Text: %s\n", b.TextVersions.URL)
}

type BillResponse struct {
	Bills []struct {
		Congress     int    `json:"congress"`
		Number       string `json:"number"`
		Title        string `json:"title"`
		UpdateDate   string `json:"updateDate"`
		LatestAction struct {
			ActionDate string `json:"actionDate"`
			Text       string `json:"text"`
		} `json:"latestAction"`
		URL string `json:"url"`
	} `json:"bills"`


}

// FetchRecentBills fetches and prints recent bills from Congress
func FetchRecentBills(apiKey string, congressNum string) error {
//	url := "https://api.congress.gov/v3/bill?congress=119&format=json"
	url := fmt.Sprintf("https://api.congress.gov/v3/bill?congress=%s&format=json", congressNum)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Api-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var response BillResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	fmt.Println("🗳️  Recent Bills in Congress:")
	for _, bill := range response.Bills {
		fmt.Printf("\n- [%d %s] %s\n  📅 %s | 🔗 %s\n",
			bill.Congress, bill.Number, bill.Title,
			bill.LatestAction.ActionDate, bill.URL)
	}

	return nil
}
func FetchVotesByBillID(apiKey, billID, congressNum string) {
    re := regexp.MustCompile(`^([a-z]+)(\d+)$`)
    matches := re.FindStringSubmatch(strings.ToLower(billID))
    if len(matches) != 3 {
        fmt.Println("Invalid bill ID format. Use something like hr2670.")
        return
    }
    billType := matches[1]
    billNumber := matches[2]

    votesURL := fmt.Sprintf("https://api.congress.gov/v3/bill/%s/%s/%s/votes?format=json", congressNum, billType, billNumber)
    req, _ := http.NewRequest("GET", votesURL, nil)
    req.Header.Set("X-Api-Key", apiKey)

    // ✅ Declare the client BEFORE using it
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Request failed:", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode == 404 {
        fmt.Printf("🚫 Congress.gov has no vote data for bill %s in congress %s\n", billID, congressNum)
        fmt.Println("💡 Try running: go run . clerkvote --year=2023 --roll=328")
        return
    }

    var data struct {
        Votes []struct {
            Chamber  string `json:"chamber"`
            RollCall string `json:"rollCallNumber"`
            Date     string `json:"date"`
            Result   string `json:"result"`
            Question string `json:"question"`
        } `json:"votes"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        fmt.Println("Decode error:", err)
        return
    }

    if len(data.Votes) == 0 {
        fmt.Println("No votes recorded for this bill.")
        return
    }

    vote := data.Votes[0]
    year := strings.Split(vote.Date, "-")[0]

    fmt.Printf("ℹ️  Found Roll Call %s (%s) in the %s on %s\n", vote.RollCall, vote.Chamber, year, vote.Date)
    FetchClerkXMLRollCall(year, vote.RollCall)
}

