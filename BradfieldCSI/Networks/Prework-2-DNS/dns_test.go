package main

import (
	"bytes"
	"testing"
)

func TestDNSMessage(t *testing.T) {

	t.Run("Encode parsed DNS Message to raw bytes", func(t *testing.T) {
		raw := sampleRawDNSMessage()
		parsed := sampleDNSMessage()
		got := parsed.encode()
		want := raw
		if !bytes.Equal(got, want) {
			t.Errorf("Expected to encode %q to %q, got %q", parsed, want, got)
		}
	})

	t.Run("Parse DNS questions from raw DNS message", func(t *testing.T) {
		raw := sampleRawDNSMessage()
		parsed := sampleDNSMessage()
		got, _ := ParseDNSQuestions(raw, parsed.Header)
		want := parsed.Questions

		if len(got) != len(want) {
			t.Errorf("Number of parsed questions != Number of raw input questions")
		}

		for i := 0; i < len(got); i++ {
			q1 := got[i]
			q2 := want[i]
			if !bytes.Equal(q1.Name, q2.Name) || q1.Type != q2.Type || q1.Class != q2.Class || q1.Namestr != q2.Namestr {
				t.Errorf("Expected parsed question %q, got %q", want, got)
			}
		}
	})

	t.Run("Parse DNS answers from DNS message", func(t *testing.T) {
		raw := sampleRawDNSMessage()
		parsed := sampleDNSMessage()
		_, offset := ParseDNSQuestions(raw, parsed.Header)
		got, _ := ParseDNSAnswers(raw, parsed.Header, offset)
		want := parsed.Answers

		if len(got) != len(want) {
			t.Errorf("Number of parsed answers != Number of raw input answers")
		}

		for i := 0; i < len(got); i++ {
			q1 := got[i]
			q2 := want[i]
			if !bytes.Equal(q1.Name, q2.Name) || q1.Type != q2.Type || q1.Class != q2.Class || q1.TTL != q2.TTL || q1.RDLength != q2.RDLength || !bytes.Equal(q1.RData, q2.RData) || q1.Namestr != q2.Namestr {
				t.Errorf("Expected parsed answer %q, got %q", want, got)
			}
		}
	})
}

/********
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
**********/
