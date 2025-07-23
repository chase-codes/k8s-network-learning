package modules

import (
	"fmt"

	"netlab/internal/tui"
	osimodel "netlab/modules/01-osi-model"
)

// ModuleInfo contains metadata about a module
type ModuleInfo struct {
	ID   string
	Name string
}

// GetModuleInfo returns the module info for a given ID
func GetModuleInfo(moduleID string) ModuleInfo {
	moduleNames := map[string]string{
		"01-osi-model":      "OSI Model Fundamentals",
		"02-tcp-ip":         "TCP/IP Stack Deep Dive",
		"03-subnetting":     "Subnetting and CIDR",
		"04-routing":        "Routing Protocols",
		"05-k8s-networking": "Kubernetes Networking",
		"06-cni":            "Container Network Interface",
		"07-service-mesh":   "Service Mesh Concepts",
	}

	name, exists := moduleNames[moduleID]
	if !exists {
		name = "Unknown Module"
	}

	return ModuleInfo{
		ID:   moduleID,
		Name: name,
	}
}

// RunModuleWithDependencyCheck runs a module with dependency checking
func RunModuleWithDependencyCheck(moduleID string) error {
	moduleInfo := GetModuleInfo(moduleID)

	// Run dependency check
	for {
		result, err := tui.RunDependencyCheck(moduleID, moduleInfo.Name)
		if err != nil {
			return fmt.Errorf("dependency check failed: %w", err)
		}

		switch result {
		case tui.CheckContinue:
			// User chose to continue, proceed with module
			return runModuleImplementation(moduleID)
		case tui.CheckAbort:
			// User chose to quit
			return nil
		case tui.CheckRecheck:
			// User chose to recheck, loop again
			continue
		default:
			return fmt.Errorf("unexpected dependency check result")
		}
	}
}

// RunModule launches a specific learning module by ID (without dependency check)
func RunModule(moduleID string) error {
	return runModuleImplementation(moduleID)
}

// runModuleImplementation contains the actual module execution logic
func runModuleImplementation(moduleID string) error {
	switch moduleID {
	case "01-osi-model", "01", "osi":
		return osimodel.Run()
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
