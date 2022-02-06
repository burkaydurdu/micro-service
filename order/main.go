package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semcov "go.opentelemetry.io/otel/semconv/v1.7.0"
)

var client http.Client

func main() {
	tp := initTracer()

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	client = http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	http.HandleFunc("/turbine/config", turbineConfigHandler)
	http.HandleFunc("/orders", createOrderHandler)
	log.Println("order-service is available at localhost:3001")
	log.Fatal(http.ListenAndServe(":3001", nil))
}

func initTracer() *sdktrace.TracerProvider {
	collectorEndpoint := os.Getenv("TRACE_COLLECTOR_ENDPOINT")
	exporter, err := zipkin.New(collectorEndpoint)

	if err != nil {
		log.Fatalf("Error creating Zipkin exporter: ", err)
	}

	fmt.Printf("Metrics will be exporter to %s", collectorEndpoint)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()), // All tracer
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semcov.SchemaURL,
			semcov.ServiceNameKey.String("order-service"))),
	)

	tmp := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(tmp)
	return tp
}

func createOrderHandler(writer http.ResponseWriter, request *http.Request) {
	ctx, span := otel.Tracer("orderTracer").Start(request.Context(), "orderSpan")

	defer span.End()

	paymentUrl := "http://localhost:8466/v1/call/payment-service/payments"

	req, _ := http.NewRequestWithContext(ctx, "POST", paymentUrl, nil)

	resp, err := client.Do(req)

	if err != nil {

	}
	// Port 8466 is sidecar's port
	// payment-service is my payment service name
	// payments is payment end point.

	//paymentRequest, err := http.NewRequest(http.MethodPost, paymentUrl, nil)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte((fmt.Sprintf("Failed to create payment request: %v", err))))
		return
	}

	//paymentResponse := getPaymentResult(paymentRequest)
	paymentResponse := getPaymentResult(resp)

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(paymentResponse))
}

func getPaymentResult(response *http.Response) string {
	defer response.Body.Close()

	bodyBytes, readError := ioutil.ReadAll(response.Body)

	if readError != nil {
		return fmt.Sprintf("Failed to read payment response: %v", readError)
	}

	return fmt.Sprintf("I send create payment request and got back %s", string(bodyBytes))
}

//func getPaymentResult(request *http.Request) string {
//	paymentResponse, err := http.DefaultClient.Do(request)
//
//	if err != nil {
//		return fmt.Sprintf("Failed to create payment: %v", err)
//	}
//
//	defer paymentResponse.Body.Close()
//	bodyBytes, readError := ioutil.ReadAll(paymentResponse.Body)
//
//	if readError != nil {
//		return fmt.Sprintf("Failed to read payment response: %v", err)
//	}
//
//	return fmt.Sprintf("I send create payment request and got back %s", string(bodyBytes))
//}

func turbineConfigHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(fmt.Sprintf(`{"serviceName": "order-service"}`)))
}
