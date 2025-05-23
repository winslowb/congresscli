package congress 

import (
	"encoding/json"
//	"os"
	"fmt"
	"net/http"
	"strings"
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
func FetchRecentBills(apiKey string) error {
	url := "https://api.congress.gov/v3/bill?congress=119&format=json"

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

func FetchVotesByBillID(apiKey, billID string) {
	parts := strings.Split(billID, "r")
	if len(parts) != 2 {
		fmt.Println("Invalid bill ID format. Example: hr2250")
		return
	}
	billType := "hr"
	billNumber := parts[1]

	url := fmt.Sprintf("https://api.congress.gov/v3/bill/119/%s/%s/votes?format=json", billType, billNumber)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Api-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching votes:", err)
		return
	}
	defer resp.Body.Close()

	var data struct {
		Votes []struct {
			Chamber     string `json:"chamber"`
			RollCall    string `json:"rollCallNumber"`
			Date        string `json:"date"`
			Result      string `json:"result"`
			Question    string `json:"question"`
			VoteURI     string `json:"url"`
		} `json:"votes"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println("Error decoding vote data:", err)
		return
	}

	if len(data.Votes) == 0 {
		fmt.Println("No votes recorded for this bill.")
		return
	}

	fmt.Printf("\n📊 Vote History for %s%s:\n", billType, billNumber)
	for _, vote := range data.Votes {
		fmt.Printf("\n• 🏛️ %s | Roll Call #%s\n", vote.Chamber, vote.RollCall)
		fmt.Printf("  📅 %s\n", vote.Date)
		fmt.Printf("  ❓ Question: %s\n", vote.Question)
		fmt.Printf("  ✅ Result: %s\n", vote.Result)
		fmt.Printf("  🔗 %s\n", vote.VoteURI)
	}
}

// Search for bills by keyword
func SearchBillsByKeyword(apiKey, keyword string) {
	url := fmt.Sprintf("https://api.congress.gov/v3/bill?query=%s&format=json", keyword)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-Api-Key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

var response BillResponse
if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
	fmt.Println("Decode error:", err)
	return
}

	for _, bill := range response.Bills {
		fmt.Printf("- [%d %s] %s\n  📅 %s | 🔗 %s\n",
			bill.Congress, bill.Number, bill.Title,
			bill.LatestAction.ActionDate, bill.URL)
	}
}
