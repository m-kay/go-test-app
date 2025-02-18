package intergration_test

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/gofiber/fiber/v2"
	"github.com/kinbiko/jsonassert"
	"github.com/m-kay/test-app/src/app"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	. "github.com/wiremock/wiremock-testcontainers-go"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type AppTest struct {
	app  *fiber.App
	Port uint
}

var testApp AppTest

func TestMain(m *testing.M) {
	ctx := context.Background()
	cont, application, err := beforeTest(ctx)
	if err != nil {
		return
	}

	m.Run()

	afterTest(application, cont, ctx)
}

func TestHello(t *testing.T) {
	//arrange
	requestUrl, err := url.Parse(fmt.Sprintf("http://localhost:%d/hello", testApp.Port))

	//act
	response, err := testApp.app.Test(&http.Request{Method: "GET", URL: requestUrl, Body: http.NoBody})

	//assert
	assert.NoError(t, err)
	if response != nil {
		assert.Equal(t, http.StatusOK, response.StatusCode)
		body, err := io.ReadAll(response.Body)
		assert.NoError(t, err)
		jsonassert.New(t).Assertf(string(body), `
		{ 
			"message":  "Hello from GRPC!"
		}`)
	}
}

func afterTest(application *app.Application, cont *WireMockContainer, ctx context.Context) {
	application.Stop()
	cont.Terminate(ctx)
}

func beforeTest(ctx context.Context) (*WireMockContainer, *app.Application, error) {
	cont, port, err := startWiremock(ctx)
	if err != nil {
		return nil, nil, err
	}

	err = os.Setenv("HELLO_SERVICE_PORT", port.Port())
	if err != nil {
		return nil, nil, err
	}

	application := startApp()
	return cont, application, nil
}

func startWiremock(ctx context.Context) (*WireMockContainer, nat.Port, error) {
	_, currentFile, _, _ := runtime.Caller(0)
	executionPath := filepath.Dir(currentFile)

	cont, err := RunContainer(ctx,
		WithImage("ghcr.io/m-kay/wiremock-grpc:latest"),
		WithMappingFile("service", executionPath+"/../../resources/wiremock/mappings/hello-service.json"),
		testcontainers.WithHostConfigModifier(func(hostConfig *container.HostConfig) {
			hostConfig.Mounts = []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: executionPath + "/../../resources/grpc",
					Target: "/home/wiremock/grpc",
				},
			}
		}),
	)
	if err != nil {
		return nil, "", fmt.Errorf("could not start container: %v", err)
	}
	port, err := cont.MappedPort(ctx, "8080")
	if err != nil {
		return nil, "", fmt.Errorf("could get container port: %v", err)
	}
	fmt.Printf("Wiremock server listening on port %d\n", port.Int())
	return cont, port, nil
}

func startApp() *app.Application {
	application := app.NewApplication()
	randomPort := uint(rand.Float32() * 9999)
	go func() {
		err := application.Start(randomPort)
		if err != nil {
			panic(err)
		}
	}()
	testApp = AppTest{app: application.Fiber, Port: randomPort}
	return application
}
