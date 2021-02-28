package interestcal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type DepositCalculation struct {
	BankName    string  `json:"bank_name,omitempty"`
	Account     string  `json:"account,omitempty"`
	AccountType string  `json:"account_type,omitempty"`
	Apy         float64 `json:"apy,omitempty"`
	Years       float64 `json:"years,omitempty"`
	Amount      float64 `json:"amount,omitempty"`
	Delta       float64 `json:"delta,omitempty"`
}

func (svc ServiceServer) writeMessage(ctx context.Context, dc DepositCalculation) {
	recordJSON, _ := json.Marshal(dc)
	log.Println("svc.AccessToken is ", svc.AccessToken)
	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("accessToken-%s", svc.AccessToken)),
		Value: recordJSON,
	}
	err := svc.KafkaWriter.WriteMessages(ctx, msg)

	if err != nil {
		log.Fatalln("could not write to kafka writer", err)
	}
}
