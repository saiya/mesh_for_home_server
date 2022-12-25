package routetable

import (
	"time"

	"github.com/saiya/mesh_for_home_server/config"
)

type Route struct {
	ExpireAt time.Time
	Dest     config.NodeID
}

func (r *Route) IsValid(now time.Time) bool {
	return r.ExpireAt.Local().After(now)
}
