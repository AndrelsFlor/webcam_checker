package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

const NOT_USED_MSG = "Câmera não está sendo usada"
const USED_MSG = "Câmera está sendo usada"
const CHAT_ID = "<ID_DA_CONVERSA>"
const URL = "https://api.telegram.org/bot<TOKEN_DO_BOT>"

func main() {

	go daemon()

	select {}

}

func sendMessage(text string) {

	url := fmt.Sprintf("%s/sendMessage", URL)

	body, err := json.Marshal(map[string]string{
		"chat_id": CHAT_ID,
		"text":    text,
	})

	if err != nil {
		fmt.Printf("Erro ao parsear json: %v\n", err)
		return
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Erro ao enviar para o telegram: %v\n", err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		body, err = io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Resposta inválida vinda da api do telegram")
			return
		}

		fmt.Printf("Status de resposta diferente de 200. Body: %v", string(body))

	}

}

func daemon() {

	isBeingUsed := false
	msg := NOT_USED_MSG
	count := 0

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		oldBeingUsed := isBeingUsed

		app := "lsmod"

		cmd := exec.Command(app)
		stdout, err := cmd.Output()

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		mods := strings.Split(string(stdout), "\n")

		for _, mod := range mods {

			modArr := strings.Fields(mod)
			if len(modArr) > 0 && modArr[0] == "uvcvideo" {
				webCamListeners := modArr[2]
				if webCamListeners == "0" {
					isBeingUsed = false
					msg = NOT_USED_MSG
				} else {
					isBeingUsed = true
					msg = USED_MSG
				}

			}
		}

		if oldBeingUsed != isBeingUsed {
			count++
		}

		if count == 1 && oldBeingUsed == isBeingUsed {
			fmt.Printf("%s\n", msg)
			count = 0
			sendMessage(msg)

		}

	}
}
