package routetable_test

import (
	"os"
	"testing"

	"github.com/saiya/mesh_for_home_server/logger"
)

func TestMain(m *testing.M) {
	logger.EnableDebugLog()
	os.Exit(m.Run())
}
