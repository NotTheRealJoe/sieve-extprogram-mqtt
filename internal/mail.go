package internal

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"strings"
)

type Email struct {
	Headers   map[string]string `json:"headers"`
	BodyParts map[string]string `json:"bodyParts"`
}

func ReadEmail(reader io.Reader) (*Email, error) {
	email := Email{
		Headers:   make(map[string]string),
		BodyParts: make(map[string]string),
	}

	msg, err := mail.ReadMessage(reader)
	if err != nil {
		return nil, err
	}

	for key, values := range msg.Header {
		// Concatenate multiple values for a single header key
		email.Headers[key] = strings.Join(values, ", ")
	}

	contentType := msg.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, fmt.Errorf("Invalid MIME type: %v", err)
	}

	if strings.HasPrefix(mediaType, "multipart/") {
		// Parse as a multipart message
		email.BodyParts, err = readMultipartEmail(msg, params)
	} else {
		email.BodyParts, err = readSinglePartEmail(msg, mediaType)
	}

	return &email, err
}

func readMultipartEmail(msg *mail.Message, mediaTypeParams map[string]string) (map[string]string, error) {
	if _, ok := mediaTypeParams["boundary"]; !ok {
		return nil, fmt.Errorf("Boundary for multipart not found in content-type parameters")
	}

	multipartReader := multipart.NewReader(msg.Body, mediaTypeParams["boundary"])

	bodyParts := make(map[string]string)

	for {
		part, err := multipartReader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error reading part: %v", err)
		}

		// Decode part based on its encoding - quoted-printable or "raw"
		var decodedPart bytes.Buffer
		if part.Header.Get("Content-Transfer-Encoding") == "quoted-printable" {
			qpReader := quotedprintable.NewReader(part)
			if _, err := io.Copy(&decodedPart, qpReader); err != nil {
				log.Fatalf("Error decoding quoted-printable part: %v", err)
			}
		} else {
			if _, err := io.Copy(&decodedPart, part); err != nil {
				log.Fatalf("Error reading part: %v", err)
			}
		}

		partMediaType, _, err := mime.ParseMediaType(part.Header.Get("Content-Type"))
		if err != nil {
			return nil, fmt.Errorf("failed to parse part's content-type header: %v", err)
		}
		bodyParts[partMediaType] = decodedPart.String()
	}

	return bodyParts, nil
}

func readSinglePartEmail(msg *mail.Message, mediaType string) (map[string]string, error) {
	var decodedBody bytes.Buffer
	if msg.Header.Get("Content-Transfer-Encoding") == "quoted-printable" {
		qpReader := quotedprintable.NewReader(msg.Body)
		if _, err := io.Copy(&decodedBody, qpReader); err != nil {
			return nil, fmt.Errorf("error decoding quoted-printable body: %v", err)
		}
	} else {
		if _, err := io.Copy(&decodedBody, msg.Body); err != nil {
			return nil, fmt.Errorf("error reading body: %v", err)
		}
	}

	return map[string]string{
		mediaType: decodedBody.String(),
	}, nil
}
