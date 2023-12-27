package gameserver

import (
	"CoGo/internal/pkg/models"
	"fmt"
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestProtobuf(t *testing.T) {
	fmt.Println("Running protobuf test")

	// protobuf struct
	profile := models.Profile{}

	mess, err := proto.Marshal(&profile)
	if err != nil {
		t.Errorf("failed to marshal proto: %v", err)
	}

	t.Log(mess)

	profile2 := models.Profile{}

	if err := proto.Unmarshal(mess, &profile2); err != nil {
		t.Errorf("failed to unmarshal proto: %v", err)
	}

	t.Log(profile2)
}
