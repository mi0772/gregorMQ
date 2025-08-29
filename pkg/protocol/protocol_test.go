package protocol

import (
	"encoding/base64"
	"regexp"
	"strings"
	"testing"
)

func encodeBracketsToBase64(input string) string {
	// Regex per trovare tutto ciò che è tra []
	re := regexp.MustCompile(`\[(.*?)\]`)

	// Funzione di sostituzione
	encoded := re.ReplaceAllStringFunc(input, func(match string) string {
		// estrai il contenuto senza le parentesi
		content := match[1 : len(match)-1]
		// codifica in base64
		encodedContent := base64.StdEncoding.EncodeToString([]byte(content))
		// ritorna con le parentesi originali
		return "[" + encodedContent + "]"
	})

	return encoded
}

func TestParseMessage(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    *Message
		shouldError bool
		errorMsg    string
	}{
		{
			name:  "Messaggio completo con header multipli",
			input: "|2{6}[client]{12}[header prova]|{10}[chiave_123]{43}[Questo è il contenuto del messaggio inviato]",
			expected: &Message{
				Magic: 0,
				Header: []MessagePart{
					{Len: 6, Content: "client"},
					{Len: 12, Content: "header prova"},
				},
				Key:     MessagePart{Len: 10, Content: "chiave_123"},
				Content: MessagePart{Len: 43, Content: "Questo è il contenuto del messaggio inviato"},
			},
			shouldError: false,
		},
		{
			name:  "Messaggio senza header",
			input: "{10}[chiave_123]{29}[Questo è il contenuto del msg]",
			expected: &Message{
				Magic:   0,
				Header:  []MessagePart{},
				Key:     MessagePart{Len: 10, Content: "chiave_123"},
				Content: MessagePart{Len: 29, Content: "Questo è il contenuto del msg"},
			},
			shouldError: false,
		},
		{
			name:  "Header con un solo elemento",
			input: "|1{4}[test]|{3}[key]{7}[content]",
			expected: &Message{
				Magic: 0,
				Header: []MessagePart{
					{Len: 4, Content: "test"},
				},
				Key:     MessagePart{Len: 3, Content: "key"},
				Content: MessagePart{Len: 7, Content: "content"},
			},
			shouldError: false,
		},
		{
			name:  "Header vuoto",
			input: "|0|{3}[abc]{5}[hello]",
			expected: &Message{
				Magic:   0,
				Header:  []MessagePart{},
				Key:     MessagePart{Len: 3, Content: "abc"},
				Content: MessagePart{Len: 5, Content: "hello"},
			},
			shouldError: false,
		},
		{
			name:  "Contenuto vuoto valido",
			input: "{4}[test]{0}[]",
			expected: &Message{
				Magic:   0,
				Header:  []MessagePart{},
				Key:     MessagePart{Len: 4, Content: "test"},
				Content: MessagePart{Len: 0, Content: ""},
			},
			shouldError: false,
		},
		// Test di errore
		{
			name:        "Header count non corrisponde",
			input:       "|2{6}[client]|{3}[key]{7}[content]",
			shouldError: true,
			errorMsg:    "numero di parti header non corrisponde",
		},
		{
			name:        "Manca Content",
			input:       "{3}[key]",
			shouldError: true,
			errorMsg:    "devono esserci esattamente 2 parti",
		},
		{
			name:        "Troppe parti",
			input:       "{3}[key]{7}[content]{5}[extra]",
			shouldError: true,
			errorMsg:    "devono esserci esattamente 2 parti",
		},
		{
			name:        "Lunghezza Key errata",
			input:       "{4}[wrong]{7}[content]",
			shouldError: true,
			errorMsg:    "errore nel parsing Key",
		},
		{
			name:        "Lunghezza Content errata",
			input:       "{3}[key]{10}[short]",
			shouldError: true,
			errorMsg:    "errore nel parsing Content",
		},
		{
			name:        "Formato completamente sbagliato",
			input:       "invalid message format",
			shouldError: true,
			errorMsg:    "devono esserci esattamente 2 parti",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input = encodeBracketsToBase64(tt.input)
			result, err := ParseMessage(tt.input)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Errorf("Expected result but got nil")
				return
			}

			// Verifica Magic
			if result.Magic != tt.expected.Magic {
				t.Errorf("Magic: expected %d, got %d", tt.expected.Magic, result.Magic)
			}

			// Verifica Header
			if len(result.Header) != len(tt.expected.Header) {
				t.Errorf("Header length: expected %d, got %d", len(tt.expected.Header), len(result.Header))
			} else {
				for i, expectedPart := range tt.expected.Header {
					if result.Header[i].Len != expectedPart.Len {
						t.Errorf("Header[%d].Len: expected %d, got %d", i, expectedPart.Len, result.Header[i].Len)
					}
					if result.Header[i].Content != expectedPart.Content {
						t.Errorf("Header[%d].Content: expected '%s', got '%s'", i, expectedPart.Content, result.Header[i].Content)
					}
				}
			}

			// Verifica Key
			if result.Key.Len != tt.expected.Key.Len {
				t.Errorf("Key.Len: expected %d, got %d", tt.expected.Key.Len, result.Key.Len)
			}
			if result.Key.Content != tt.expected.Key.Content {
				t.Errorf("Key.Content: expected '%s', got '%s'", tt.expected.Key.Content, result.Key.Content)
			}

			// Verifica Content
			if result.Content.Len != tt.expected.Content.Len {
				t.Errorf("Content.Len: expected %d, got %d", tt.expected.Content.Len, result.Content.Len)
			}
			if result.Content.Content != tt.expected.Content.Content {
				t.Errorf("Content.Content: expected '%s', got '%s'", tt.expected.Content.Content, result.Content.Content)
			}
		})
	}
}

func TestParseMessagePart(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    *MessagePart
		shouldError bool
	}{
		{
			name:     "Parte valida normale",
			input:    "{10}[chiave_123]",
			expected: &MessagePart{Len: 10, Content: "chiave_123"},
		},
		{
			name:     "Parte vuota valida",
			input:    "{0}[]",
			expected: &MessagePart{Len: 0, Content: ""},
		},
		{
			name:     "Contenuto con spazi",
			input:    "{12}[hello world!]",
			expected: &MessagePart{Len: 12, Content: "hello world!"},
		},
		{
			name:        "Lunghezza non corrisponde",
			input:       "{15}[short]",
			shouldError: true,
		},
		{
			name:        "Formato non valido",
			input:       "invalid",
			shouldError: true,
		},
		{
			name:        "Lunghezza non numerica",
			input:       "{abc}[test]",
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.input = encodeBracketsToBase64(tt.input)
			result, err := parseMessagePart(tt.input)

			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result.Len != tt.expected.Len {
				t.Errorf("Len: expected %d, got %d", tt.expected.Len, result.Len)
			}
			if result.Content != tt.expected.Content {
				t.Errorf("Content: expected '%s', got '%s'", tt.expected.Content, result.Content)
			}
		})
	}
}

// Benchmark per testare le performance
func BenchmarkParseMessage(b *testing.B) {
	input := "|2{6}[client]{12}[header prova]|{10}[chiave_123]{43}[Questo è il contenuto del messaggio inviato]"
	input = encodeBracketsToBase64(input)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ParseMessage(input)
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkParseMessagePart(b *testing.B) {
	input := "{10}[chiave_123]"
	input = encodeBracketsToBase64(input)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := parseMessagePart(input)
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}
