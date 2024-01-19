package helpers

import (
	"strings"

	"github.com/yash492/statusy/pkg/types"
)

func ConvertComponentsToStr(components []types.ComponentsWithNameAndID) string {
	var sb strings.Builder
	for i, component := range components {
		sb.WriteString(component.Name)
		if i < len(components)-1 {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}

// incident_update_state -> discord -> #45545
var ChatopsMsgColor map[string]map[string]string
