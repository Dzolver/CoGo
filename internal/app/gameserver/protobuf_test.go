package gameserver

import (
	"CoGo/internal/pkg/models"
	"fmt"
	"testing"
)

func TestProtobuf(t *testing.T) {
	fmt.Println("Running protobuf test")

	protoClient := models.Client{}
	protoClient.GetBroadcastAddr()
}
