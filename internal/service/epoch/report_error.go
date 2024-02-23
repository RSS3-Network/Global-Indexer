package epoch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"mime/multipart"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
)

type SlackMessageElement struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type SlackMessageBlock struct {
	Type     string                `json:"type"`
	Text     *SlackMessageElement  `json:"text,omitempty"`
	Elements []SlackMessageElement `json:"elements,omitempty"`
}

type SlackMessageStructure struct {
	Channel string              `json:"channel"`
	Text    string              `json:"text"` // For fallback and notification
	Blocks  []SlackMessageBlock `json:"blocks"`
}

func (s *Server) ReportFailedTransactionToSlack(txErr error, txHash string, txFunc string, users []common.Address, amount []*big.Int) {
	log.Println("===================== Transaction Error =====================")
	log.Printf("%s transaction error! Please check information below:", txFunc)
	log.Printf("Transaction Hash: %s", txHash)
	log.Printf("Error: %v", txErr)
	log.Println("--------------- Users List (address, amount) ---------------")

	for i, user := range users {
		log.Printf("%s,%s", user.Hex(), amount[i].String())
	}

	log.Println("------------------------------------------------------------")
	log.Printf("Preparing to send information to slack...")

	defer log.Println("============================================================")

	// Send base information
	var errReason string

	if errors.Is(txErr, errors.New("transaction failed")) {
		errReason = "Failed"
	} else if errors.Is(txErr, context.DeadlineExceeded) {
		errReason = "Timeout (/!\\ doesn't mean it failed)"
	} else {
		// Unknown error
		errReason = txErr.Error()
	}

	summary := fmt.Sprintf("⚠ %s transaction error", txFunc)
	message := SlackMessageStructure{
		Channel: s.slackNotificationChannel,
		Text:    summary,
		Blocks: []SlackMessageBlock{
			{
				Type: "header",
				Text: &SlackMessageElement{
					"plain_text",
					summary,
				},
			},
			{
				Type: "section",
				Text: &SlackMessageElement{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Error Message*: %s", errReason),
				},
			},
			{
				Type: "context",
				Elements: []SlackMessageElement{
					{
						Type: "mrkdwn",
						Text: fmt.Sprintf("*Txn*: <%s%s|%s>", s.slackNotificationBlockchainScan, txHash, txHash),
					},
				},
			},
		},
	}

	messageBytes, err := json.Marshal(&message)

	if err != nil {
		log.Printf("Failed to parse into json with error: %v", err)
		return
	}

	msgReq, err := http.NewRequestWithContext(context.Background(), "POST", "https://slack.com/api/chat.postMessage", bytes.NewReader(messageBytes))

	if err != nil {
		log.Printf("Failed to prepare message request with error: %v", err)
		return
	}

	msgReq.Header.Set("Content-Type", "application/json")
	msgReq.Header.Set("Authorization", "Bearer "+s.slackNotificationBotToken)

	res, err := (&http.Client{}).Do(msgReq)

	if err != nil {
		log.Printf("Failed to send error log to Slack with error: %v", err)
		return
	}

	res.Body.Close()

	log.Printf("Notification message sent, preparing users list...")

	// Upload failed users list as csv file
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)
	channelsWriter, err := bodyWriter.CreateFormField("channels")

	if err != nil {
		log.Printf("Failed to prepare channels form with error: %v", err)
		return
	}

	_, err = channelsWriter.Write([]byte(s.slackNotificationChannel))

	if err != nil {
		log.Printf("Failed to write channels field with error: %v", err)
		return
	}

	fileWriter, err := bodyWriter.CreateFormFile(
		"file",
		fmt.Sprintf("user-and-amount-%s.csv", txHash),
	)

	if err != nil {
		log.Printf("Failed to prepare file upload form with error: %v", err)
		return
	}

	_, err = fileWriter.Write([]byte("address,amount\n"))

	if err != nil {
		log.Printf("Failed to write csv title into file with error: %v", err)
		return
	}

	for i, user := range users {
		_, err = fileWriter.Write([]byte(fmt.Sprintf("%s,%s\n", user.Hex(), amount[i].String())))
		if err != nil {
			log.Printf("Failed to write user line (%s,%s) into file with error: %v", user.Hex(), amount[i].String(), err)
		}
	}

	// Finish write, upload to slack
	bodyWriter.Close()

	fileReq, err := http.NewRequestWithContext(context.Background(), "POST", "https://slack.com/api/files.upload", bodyBuffer)

	if err != nil {
		log.Printf("Failed to prepare file request with error: %v", err)
		return
	}

	fileReq.Header.Set("Content-Type", bodyWriter.FormDataContentType())
	fileReq.Header.Set("Authorization", "Bearer "+s.slackNotificationBotToken)

	res, err = (&http.Client{}).Do(fileReq)

	if err != nil {
		log.Printf("Failed to send error log to Slack with error: %v", err)
		return
	}

	res.Body.Close()

	log.Printf("Users list file uploaded.")
}
