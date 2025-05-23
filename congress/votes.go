package congress

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type RollCallVote struct {
	XMLName       xml.Name     `xml:"rollcall-vote"`
	VoteMetadata  VoteMeta     `xml:"vote-metadata"`
	RecordedVotes []Member     `xml:"recorded-vote"`
}

type VoteMeta struct {
	VoteDate string `xml:"action-date"`
	LegisNum string `xml:"legis-num"`
	Question string `xml:"vote-question"`
	Result   string `xml:"vote-result"`
	VoteDesc string `xml:"vote-desc"` // ✅ moved here
	PartyTotals 	[]PartyTotal 		`xml:"vote-totals>totals-by-party"`

}

type Legislation struct {
	Type   string `xml:"type"`
	Number string `xml:"number"`
}

type Totals struct {
	Yeas      int `xml:"yea"`
	Nays      int `xml:"nay"`
	Present   int `xml:"present"`
	NotVoting int `xml:"not-voting"`
}

type Member struct {
	Name  string `xml:"legislator,attr"`
	Party string `xml:"party,attr"`
	State string `xml:"state,attr"`
	Vote  string `xml:",chardata"`
}

type PartyTotal struct {
	Party     string `xml:"party"`
	Yeas      int    `xml:"yea-total"`
	Nays      int    `xml:"nay-total"`
	Present   int    `xml:"present-total"`
	NotVoting int    `xml:"not-voting-total"`
}


func FetchClerkXMLRollCall(year string, rollNumber string) {
	url := fmt.Sprintf("https://clerk.house.gov/evs/%s/roll%s.xml", year, rollNumber)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}


if resp.StatusCode != http.StatusOK {
	fmt.Printf("HTTP error: %d %s\n", resp.StatusCode, resp.Status)
	return
}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

fmt.Println("===== RAW XML PREVIEW =====")
fmt.Println(string(body))
fmt.Println("===========================")


var vote RollCallVote
if err := xml.Unmarshal(body, &vote); err != nil {
	fmt.Println("XML parse error:", err)
	return
	}

// fmt.Printf("DEBUG: Found %d party total entries\n", len(vote.VoteTotals.PartyTotals))


fmt.Printf("DEBUG: Found %d party total entries\n", len(vote.VoteMetadata.PartyTotals))

fmt.Println("\n🧮 Vote Totals by Party:")
for _, pt := range vote.VoteMetadata.PartyTotals {
	fmt.Printf("• %s: 🟢 %d | 🔴 %d | ⚪ %d | ❌ %d\n",
		pt.Party, pt.Yeas, pt.Nays, pt.Present, pt.NotVoting)
}

 	// Header and metadata
 	fmt.Printf("\n📜 Roll Call Vote %s (%s)\n", rollNumber, year)
 	fmt.Printf("🗓️ Date: %s\n", vote.VoteMetadata.VoteDate)
 	fmt.Printf("📜 Bill: %s — %s\n", vote.VoteMetadata.LegisNum, vote.VoteMetadata.VoteDesc)
 	fmt.Printf("❓ Question: %s\n", vote.VoteMetadata.Question)
 	fmt.Printf("✅ Result: %s\n", vote.VoteMetadata.Result)
 	// (Raw XML preview removed for clarity)
//fmt.Println("\n🧮 Vote Totals by Party:")
//for _, pt := range vote.VoteTotals.PartyTotals {
//	fmt.Printf("• %s: 🟢 %d | 🔴 %d | ⚪ %d | ❌ %d\n",
//		pt.Party, pt.Yeas, pt.Nays, pt.Present, pt.NotVoting)
//}

	fmt.Println("\n🧑‍🤝‍🧑 Sample Votes:")
	for i, m := range vote.RecordedVotes {
		fmt.Printf("- %s (%s-%s): %s\n", m.Name, m.Party, m.State, m.Vote)
		if i == 9 {
			fmt.Println("... (more omitted)")
			break
		}
	}
}
