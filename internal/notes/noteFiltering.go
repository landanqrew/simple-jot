package notes

import (
	"fmt"
	"log"
	"time"
)

// FilterNotesByDate filters notes by date. date format for input is YYYY-MM-DD
func FilterNotesByDate(notes []Note, startDate string, endDate string) []Note {
	if startDate == "" {
		if endDate == "" {
			return notes
		}
		
		ed, err := time.Parse(time.DateOnly, endDate)
		if err != nil {
			log.Fatal("failed to parse end date: " + err.Error() + " use format YYYY-MM-DD")
		}
		filteredNotes := make([]Note, 0)
		for _, note := range notes {
			nd, err := time.Parse(time.DateTime, note.CreatedAt)
			if err != nil {
				fmt.Println("failed to parse note date: " + err.Error())
				// include regardless
				filteredNotes = append(filteredNotes, note)
			}
			if nd.Before(ed) || nd.Equal(ed) {
				filteredNotes = append(filteredNotes, note)
			}
		}
		return filteredNotes
	} else {
		if endDate == "" {
			sd, err := time.Parse(time.DateOnly, startDate)
			if err != nil {
				log.Fatal("failed to parse start date: " + err.Error() + " use format YYYY-MM-DD")
			}
			filteredNotes := make([]Note, 0)
			for _, note := range notes {
				nd, err := time.Parse(time.DateTime, note.CreatedAt)
				if err != nil {
					fmt.Println("failed to parse note date: " + err.Error())
					// include regardless
					filteredNotes = append(filteredNotes, note)
				}
				if nd.After(sd) || nd.Equal(sd)  {
					filteredNotes = append(filteredNotes, note)
				}
			}
			return filteredNotes
		} else {
			sd, err := time.Parse(time.DateOnly, startDate)
			if err != nil {
				log.Fatal("failed to parse start date: " + err.Error() + " use format YYYY-MM-DD")
			}
			// fmt.Println("sd: " + sd.Format(time.DateTime))
			ed, err := time.Parse(time.DateOnly, endDate)
			if err != nil {
				log.Fatal("failed to parse end date: " + err.Error() + " use format YYYY-MM-DD")
			}
			// fmt.Println("ed: " + ed.Format(time.DateTime))
			filteredNotes := make([]Note, 0)
			for _, note := range notes {
				nd, err := time.Parse(time.DateTime, note.CreatedAt)
				if err != nil {
					fmt.Println("failed to parse note date: " + err.Error())
					// include regardless
					filteredNotes = append(filteredNotes, note)
				}
				// fmt.Println("nd: " + nd.Format(time.DateTime))
				if (nd.After(sd) || nd.Equal(sd)) && (nd.Before(ed) || nd.Equal(ed)) {
					filteredNotes = append(filteredNotes, note)
				}
			}
			return filteredNotes
		}
	}
}

