package readme_test

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/liveoaklabs/readme-api-go-client/tests/testdata"
	"github.com/stretchr/testify/assert"
)

func Test_OutboundIPs_Get(t *testing.T) {
	t.Run("when API responds with 200", func(t *testing.T) {
		// Arrange
		expect := testdata.OutboundIPs
		gock.New(TestClient.APIURL).
			Get("/outbound-ips").
			Reply(200).
			JSON(expect)
		defer gock.Off()

		// Act
		got, _, err := TestClient.OutboundIP.Get()

		// Assert
		assert.NoError(t, err, "it does not return an error")
		assert.Equal(t, expect, got, "it returns expected data struct")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})

	t.Run("when API response cannot be parsed", func(t *testing.T) {
		// Arrange
		gock.New(TestClient.APIURL).
			Get("/outbound-ips").
			Reply(200).
			JSON(`[{"invalid":invalid"}]`)
		defer gock.Off()

		// Act
		_, _, err := TestClient.OutboundIP.Get()

		// Assert
		assert.ErrorContains(t, err,
			"unable to parse API response: invalid character",
			"it returns the expected error")
		assert.True(t, gock.IsDone(), "it makes the expected API call")
	})
}
