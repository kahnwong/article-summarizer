package cmd

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Strubbl/wallabago/v9"
	"github.com/gofiber/fiber/v3"
)

type failingWallabagClient struct{}

func (failingWallabagClient) GetEntries() ([]wallabago.Item, error) {
	return nil, errors.New("unavailable")
}

func (failingWallabagClient) MarkEntryAsRead(int) error {
	return nil
}

func TestRootControllerReturnsWallabagError(t *testing.T) {
	previousClient := wallabagClient
	wallabagClient = failingWallabagClient{}
	t.Cleanup(func() {
		wallabagClient = previousClient
	})

	app := fiber.New()
	app.Get("/", rootController)

	response, err := app.Test(httptest.NewRequest(http.MethodGet, "/", nil))
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("io.ReadAll() error = %v", err)
	}
	if response.StatusCode != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", response.StatusCode, http.StatusInternalServerError)
	}
	if got := string(body); !strings.Contains(got, "cannot obtain articles from Wallabag: unavailable") {
		t.Errorf("body = %q, want Wallabag error", got)
	}
}
