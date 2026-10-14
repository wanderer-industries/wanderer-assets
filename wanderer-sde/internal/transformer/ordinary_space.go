package transformer

import (
	"fmt"
	"math"
	"strings"

	"github.com/guarzo/wanderer-sde/internal/models"
)

// normalizeWormholeClasses supplements the sparse SDE assignments with ordinary
// space classifications, without overriding explicit or inherited classes.
func (t *Transformer) normalizeWormholeClasses(classes []models.WormholeClassLocation, systems []models.SolarSystem) ([]models.WormholeClassLocation, error) {
	result := make([]models.WormholeClassLocation, len(classes))
	copy(result, classes)
	assigned := make(map[int64]struct{}, len(classes))
	for _, class := range classes {
		assigned[class.LocationID] = struct{}{}
	}

	for _, system := range systems {
		_, systemAssigned := assigned[system.SolarSystemID]
		_, constellationAssigned := assigned[system.ConstellationID]
		_, regionAssigned := assigned[system.RegionID]
		if systemAssigned || constellationAssigned || regionAssigned {
			continue
		}

		// securityClass is ordinary-space evidence, not the security band itself.
		// Special systems without this evidence must remain unclassified.
		if strings.TrimSpace(system.SecurityClass) == "" {
			continue
		}
		security := system.Security
		if math.IsNaN(security) || math.IsInf(security, 0) || security < -1 || security > 1 {
			return nil, fmt.Errorf("solar system %d has invalid security %v: expected a finite value between -1 and 1", system.SolarSystemID, security)
		}

		classID := int64(8)
		if security <= 0 {
			classID = 9
		} else if GetTrueSecurity(security) >= 0.5 {
			classID = 7
		}
		result = append(result, models.WormholeClassLocation{
			LocationID:      system.SolarSystemID,
			WormholeClassID: classID,
		})
		assigned[system.SolarSystemID] = struct{}{}
	}

	return t.sortWormholeClasses(result), nil
}
