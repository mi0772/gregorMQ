package protocol

import (
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

type StatusResponse uint8

const (
	OK StatusResponse = iota
	KO
)

type MessagePart struct {
	Len     int
	Content string
}

type Message struct {
	Magic   uint64
	Header  []MessagePart
	Key     MessagePart
	Content MessagePart
}

type Response struct {
	status  StatusResponse
	message string
}

// ParseMessage parsa una stringa completa nel formato del messaggio
func ParseMessage(input string) (*Message, error) {
	var headerParts []MessagePart
	remaining := input

	// Controlla se c'è un header opzionale |numero{...}[...]...|
	headerRegex := regexp.MustCompile(`^\|(\d+)(.+?)\|(.*)$`)
	headerMatches := headerRegex.FindStringSubmatch(input)

	if headerMatches != nil {
		// C'è un header
		headerCount, err := strconv.Atoi(headerMatches[1])
		if err != nil {
			return nil, fmt.Errorf("numero header non valido: %v", err)
		}

		headerContent := headerMatches[2]
		remaining = headerMatches[3]

		// Parsa le parti dell'header
		partRegex := regexp.MustCompile(`\{\d+\}\[[^\]]*\]`)
		partMatches := partRegex.FindAllString(headerContent, -1)

		if len(partMatches) != headerCount {
			return nil, fmt.Errorf("numero di parti header non corrisponde: dichiarate %d, trovate %d", headerCount, len(partMatches))
		}

		// Parsa ogni parte dell'header
		for _, partStr := range partMatches {
			part, err := parseMessagePart(partStr)
			if err != nil {
				return nil, fmt.Errorf("errore nel parsing header part: %v", err)
			}
			headerParts = append(headerParts, *part)
		}
	}

	// Parsa Key e Content (obbligatori)
	partRegex := regexp.MustCompile(`\{\d+\}\[[^\]]*\]`)
	remainingParts := partRegex.FindAllString(remaining, -1)

	if len(remainingParts) != 2 {
		return nil, fmt.Errorf("devono esserci esattamente 2 parti (Key e Content), trovate %d", len(remainingParts))
	}

	// Parsa Key
	keyPart, err := parseMessagePart(remainingParts[0])
	if err != nil {
		return nil, fmt.Errorf("errore nel parsing Key: %v", err)
	}

	// Parsa Content
	contentPart, err := parseMessagePart(remainingParts[1])
	if err != nil {
		return nil, fmt.Errorf("errore nel parsing Content: %v", err)
	}

	return &Message{
		Magic:   0, // Magic non presente nell'input, impostato a 0
		Header:  headerParts,
		Key:     *keyPart,
		Content: *contentPart,
	}, nil
}

// parseMessagePart parsa una singola parte nel formato {len}[content]
func parseMessagePart(input string) (*MessagePart, error) {
	re := regexp.MustCompile(`^\{(\d+)\}\[([^\]]*)\]`)
	matches := re.FindStringSubmatch(input)

	if matches == nil {
		return nil, errors.New("formato MessagePart non valido: deve essere {numero}[contenuto]")
	}

	declaredLen, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("lunghezza non valida: %v", err)
	}

	contentB, err := base64.StdEncoding.DecodeString(matches[2])
	content := string(contentB)
	if err != nil {
		return nil, fmt.Errorf("impossibile decodificare il content : %v", err)
	}
	actualLen := utf8.RuneCountInString(string(content))

	if declaredLen != actualLen {
		return nil, fmt.Errorf("lunghezza non corrisponde: dichiarata %d, reale %d", declaredLen, actualLen)
	}

	return &MessagePart{
		Len:     declaredLen,
		Content: content,
	}, nil
}

// String restituisce una rappresentazione stringa del Message
func (m *Message) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Magic: %d\n", m.Magic))

	if len(m.Header) > 0 {
		sb.WriteString("Header:\n")
		for i, part := range m.Header {
			sb.WriteString(fmt.Sprintf("  [%d] Len: %d, Content: '%s'\n", i, part.Len, part.Content))
		}
	} else {
		sb.WriteString("Header: nessuno\n")
	}

	sb.WriteString(fmt.Sprintf("Key: Len: %d, Content: '%s'\n", m.Key.Len, m.Key.Content))
	sb.WriteString(fmt.Sprintf("Content: Len: %d, Content: '%s'\n", m.Content.Len, m.Content.Content))

	return sb.String()
}
