package modules

import (
	"fmt"

	"netlab/internal/modules/osimodel"
)

// RunModule launches a specific learning module by ID
func RunModule(moduleID string) error {
	switch moduleID {
	case "01-osi-model":
		return osimodel.RunEnhanced()
	case "02-tcp-ip":
		return fmt.Errorf("module %s is not implemented yet", moduleID)
	case "03-subnetting":
		return fmt.Errorf("module %s is not implemented yet", moduleID)
	case "04-routing":
		return fmt.Errorf("module %s is not implemented yet", moduleID)
	case "05-k8s-networking":
		return fmt.Errorf("module %s is not implemented yet", moduleID)
	case "06-cni":
		return fmt.Errorf("module %s is not implemented yet", moduleID)
	case "07-service-mesh":
		return fmt.Errorf("module %s is not implemented yet", moduleID)
	default:
		return fmt.Errorf("unknown module: %s", moduleID)
	}
}

// ListModules returns all available modules
func ListModules() []string {
	return []string{
		"01-osi-model",
		"02-tcp-ip",
		"03-subnetting",
		"04-routing",
		"05-k8s-networking",
		"06-cni",
		"07-service-mesh",
	}
}
