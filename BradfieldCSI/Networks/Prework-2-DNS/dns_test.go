package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestDNSMessage(t *testing.T) {
	t.Run("No questions exist error", func(t *testing.T) {
		_, err := NewDNSQuery(1, true, []DNSQuestion{})
		want := fmt.Errorf("No questions exist")
		if err == nil {
			t.Errorf("Expected %q error, got none", want)
		}
	})

	t.Run("Encode query to bytes", func(t *testing.T) {
		queryStruct, queryRaw := TestQueryData()
		got := queryStruct.encode()
		want := queryRaw
		if !bytes.Equal(got, want) {
			t.Errorf("Expected to encode %q to %q, got %q", queryStruct, want, got)
		}
	})

	// t.Run("Decode DNS Message into bytes", func(t *testing.T) {
	// 	q, _ := TestQueryData()
	// 	got := q.execute()
	// 	want := []byte{0x01}
	// 	if !bytes.Equal(got, want) {
	// 		t.Errorf("Expected DNS query response %q, got %q", want, got)
	// 	}
	// })

	// t.Run("Decode bytes into DNS Message", func(t *testing.T) {
	// 	queryStruct, queryRaw := TestQueryData()
	// 	got := queryStruct.encode()
	// 	want := queryRaw
	// 	if !bytes.Equal(got, want) {
	// 		t.Errorf("Expected to encode %q to %q, got %q", queryStruct, want, got)
	// 	}
	// })

	// q.Header.QCount < len(q.Questions) {
	// 	t.Errorf
	// }
	// t.Errorf()
}

func TestDNSResponse(t *testing.T) {
	t.Run("Parse DNS questions from raw DNS message", func(t *testing.T) {
		parsedDNSQuery, _ := TestQueryData()
		got, _ := ParseDNSQuestions(parsedDNSQuery.encode(), parsedDNSQuery.Header)
		want := parsedDNSQuery.Questions

		if len(got) != len(want) {
			t.Errorf("Number of parsed questions != Number of raw input questions")
		}

		for i := 0; i < len(got); i++ {
			q1 := got[i]
			q2 := want[i]
			if !bytes.Equal(q1.Name, q2.Name) || q1.Type != q2.Type || q1.Class != q2.Class {
				t.Errorf("Expected parsed question %q, got %q", want, got)
			}
		}
	})

	t.Run("Parse DNS answers from DNS message", func(t *testing.T) {
		parsedDNSResponse := sampleDNSResponse()
		rawDNSResponse := sampleRawDNSResponse()
		_, offset := ParseDNSQuestions(rawDNSResponse, parsedDNSResponse.Header)
		got, _ := ParseDNSAnswers(rawDNSResponse, offset)
		want := parsedDNSResponse.Answers

		if len(got) != len(want) {
			t.Errorf("Number of parsed answers != Number of raw input answers")
		}

		for i := 0; i < len(got); i++ {
			q1 := got[i]
			q2 := want[i]
			if !bytes.Equal(q1.Name, q2.Name) || q1.Type != q2.Type || q1.Class != q2.Class || q1.TTL != q2.TTL || q1.RDLength != q2.RDLength || !bytes.Equal(q1.RData, q2.RData) {
				t.Errorf("Expected parsed answer %q, got %q", want, got)

				fmt.Printf("\nName %x", q1.Name)
				fmt.Printf("\nType %x", q1.Type)
				fmt.Printf("\nClass %x", q1.Class)
				fmt.Printf("\nTTL %x", q1.TTL)
				fmt.Printf("\nRDLength %x", q1.RDLength)
				fmt.Printf("\nRData %x", q1.RData)

				fmt.Printf("\nName %x", q2.Name)
				fmt.Printf("\nType %x", q2.Type)
				fmt.Printf("\nClass %x", q2.Class)
				fmt.Printf("\nTTL %x", q2.TTL)
				fmt.Printf("\nRDLength %x", q2.RDLength)
				fmt.Printf("\nRData %x", q2.RData)
			}
		}
	})
}
