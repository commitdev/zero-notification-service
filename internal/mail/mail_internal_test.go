package mail

import (
	"testing"

	"github.com/commitdev/zero-notification-service/internal/server"
	"github.com/stretchr/testify/assert"
)

func TestConvertAddresses(t *testing.T) {

	addresses := []server.EmailRecipient{
		server.EmailRecipient{Name: "User1", Address: "user1@address.com"},
	}
	convertedAddresses := convertAddresses(addresses)

	assert.Len(t, convertedAddresses, 1, "Returned list should be 1 element")
	assert.Equal(t, convertedAddresses[0].Name, "User1", "Returned name should match")
	assert.Equal(t, convertedAddresses[0].Address, "user1@address.com", "Returned address should match")
}
