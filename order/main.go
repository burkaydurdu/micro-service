package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/turbine/config", turbineConfigHandler)
	http.HandleFunc("/orders", createOrderHandler)
	log.Println("order-service is available at localhost:3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
}

func createOrderHandler(writer http.ResponseWriter, request *http.Request) {
	// Port 8466 is sidecar's port
	// payment-service is my payment service name
	// payments is payment end point.
	paymentUrl := "http://localhost:8466/v1/call/payment-service/payments"
	paymentRequest, err := http.NewRequest(http.MethodPost, paymentUrl, nil)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte((fmt.Sprintf("Failed to create payment request: %v", err))))
		return
	}

	paymentResponse := getPaymentResult(paymentRequest)

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(paymentResponse))
}

func getPaymentResult(request *http.Request) string {
	paymentResponse, err := http.DefaultClient.Do(request)

	if err != nil {
		return fmt.Sprintf("Failed to create payment: %v", err)
	}

	defer paymentResponse.Body.Close()
	bodyBytes, readError := ioutil.ReadAll(paymentResponse.Body)

	if readError != nil {
		return fmt.Sprintf("Failed to read payment response: %v", err)
	}

	return fmt.Sprintf("I send create payment request and got back %s", string(bodyBytes))
}

func turbineConfigHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(fmt.Sprintf(`{"serviceName": "order-service"}`)))
}
