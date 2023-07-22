package chatservices

import (
	"encoding/json"
	"log"
)

func (message *Message) Encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
		return nil
	}

	return json
}
