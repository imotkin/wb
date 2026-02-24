package main

import (
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"
)

var statuses = []string{"Готово", "В работе", "Не будет сделано"}

type Ticket struct {
	Ticket string
	User   string
	Status string
	Date   time.Time
}

func GetTasks(text string, user *string, status *string) []Ticket {
	var tickets []Ticket

	lines := strings.Split(text, "\n")

	var filter func(user, state string) bool

	switch {
	case user == nil && status == nil:
		filter = func(_, _ string) bool {
			return true
		}
	case user == nil && status != nil:
		filter = func(_, s string) bool {
			return s == *status
		}
	case user != nil && status == nil:
		filter = func(u, _ string) bool {
			return u == *user
		}
	default:
		filter = func(u, s string) bool {
			return u == *user && s == *status
		}
	}

	for _, line := range lines {
		words := strings.Split(line, "_")

		if len(words) != 4 {
			continue
		}

		if !strings.HasPrefix(words[0], "TICKET-") {
			continue
		}

		if filter(words[1], words[2]) {
			date, err := time.Parse("2006-01-02", words[3])
			if err != nil {
				continue
			}

			if !slices.Contains(statuses, words[2]) {
				continue
			}

			tickets = append(tickets, Ticket{
				Ticket: words[0],
				User:   words[1],
				Status: words[2],
				Date:   date,
			})
		}
	}

	return tickets
}

func String(s string) *string {
	return &s
}

func TestGetTasks(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		user    *string
		status  *string
		tickets []Ticket
	}{
		{
			input: []string{
				"TICKET-12345_Паша Попов_Готово_2024-01-01",
				"TICKET-12346_Иван Иванов_В работе_2024-01-02",
				"TICKET-12347_Анна Смирнова_Не будет сделано_2024-01-03",
				"TICKET-12348_Паша Попов_В работе_2024-01-04",
			},
			user:   String("Паша Попов"),
			status: String("Готово"),
			tickets: []Ticket{
				{
					Ticket: "TICKET-12345",
					User:   "Паша Попов",
					Status: "Готово",
					Date:   time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			text := strings.Join(tt.input, "\n")

			tickets := GetTasks(text, tt.user, tt.status)
			if !reflect.DeepEqual(tt.tickets, tickets) {
				t.Errorf("got: %v, want: %v\n", tickets, tt.tickets)
			}
		})
	}
}
