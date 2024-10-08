package manifest

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Dependency map: resources dependent on others
var dependencyMap = map[string][]string{
	"Observation":        {"Encounter", "Patient", "Practitioner"},
	"Condition":          {"Encounter", "Patient", "Practitioner"},
	"Medication":         {"Patient"},
	"MedicationRequest":  {"Patient", "Practitioner", "Encounter", "Medication"},
	"DiagnosticReport":   {"Encounter", "Patient", "Practitioner", "Observation"},
	"Immunization":       {"Patient", "Practitioner", "Encounter"},
	"Procedure":          {"Encounter", "Patient", "Practitioner"},
	"AllergyIntolerance": {"Patient", "Practitioner"},

	"Encounter":         {"Patient", "Practitioner", "Location", "Organization"},
	"Patient":           {},
	"Practitioner":      {},
	"Organization":      {},
	"Location":          {"Organization"},
	"HealthcareService": {"Organization", "Location"},
	"Coverage":          {"Patient", "Organization"},

	"Claim":                {"Patient", "Practitioner", "Organization", "Coverage", "Encounter"},
	"ExplanationOfBenefit": {"Patient", "Practitioner", "Organization", "Coverage", "Claim"},

	"Appointment": {"Patient", "Practitioner", "Encounter", "Location"},
	"Schedule":    {"Practitioner", "Location"},
	"Slot":        {"Schedule", "Location"},

	"MedicationAdministration": {"Encounter", "Patient", "Practitioner", "Medication"},
	"MedicationStatement":      {"Patient", "Medication"},

	"CarePlan": {"Patient", "Practitioner", "Encounter"},
	"CareTeam": {"Patient", "Practitioner"},
	"Goal":     {"Patient", "CarePlan"},

	"ServiceRequest":    {"Patient", "Practitioner", "Encounter"},
	"Specimen":          {"Patient", "Practitioner", "Encounter", "ServiceRequest"},
	"ImagingStudy":      {"Patient", "Practitioner", "Encounter", "Specimen"},
	"DocumentReference": {"Patient", "Practitioner", "Encounter"},
}

type Manifest struct {
	ResourceType string      `json:"resourceType"`
	Parameter    []Parameter `json:"parameter"`
}

type Parameter struct {
	Name  string `json:"name"`
	Part  []Part `json:"part,omitempty"`
	Value string `json:"valueString,omitempty"`
}

type Part struct {
	Name  string `json:"name"`
	Value string `json:"valueString"`
}

// CreateManifest generates FHIR manifest JSON structures, split by dependency level
func CreateManifest(files []string, urlBase string) []Manifest {
	resourceFiles := groupFilesByResourceType(files)

	// Track which resources have been added and are "resolved"
	resolved := make(map[string]bool)
	manifests := []Manifest{}

	for len(resourceFiles) > 0 {
		manifest := Manifest{
			ResourceType: "Parameters",
			Parameter: []Parameter{
				{
					Name:  "inputFormat",
					Value: "application/fhir+ndjson",
				},
			},
		}

		var resolvedInCurrentItr []string

		for resourceType, files := range resourceFiles {
			if canAddResource(resourceType, resolved) {
				for _, file := range files {
					fileURL := fmt.Sprintf("%s/%s", urlBase, filepath.Base(file))

					inputPart := Parameter{
						Name: "input",
						Part: []Part{
							{Name: "type", Value: resourceType},
							{Name: "url", Value: fileURL},
						},
					}

					manifest.Parameter = append(manifest.Parameter, inputPart)
				}

				resolvedInCurrentItr = append(resolvedInCurrentItr, resourceType)
				delete(resourceFiles, resourceType)
			}
		}

		for _, resourceType := range resolvedInCurrentItr {
			resolved[resourceType] = true
		}
		manifests = append(manifests, manifest)
	}

	return manifests
}

// Helper function to check if a resource can be added based on its dependencies
func canAddResource(resourceType string, resolved map[string]bool) bool {
	// Check if all dependencies for this resource are resolved
	for _, dep := range dependencyMap[resourceType] {
		if !resolved[dep] {
			return false
		}
	}
	return true
}

// Helper function to infer resource type from file name
func inferResourceType(filePath string) string {
	base := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	parts := strings.Split(base, ".")
	return parts[0]
}

// Group files by their inferred resource type
func groupFilesByResourceType(files []string) map[string][]string {
	resourceFiles := make(map[string][]string)
	for _, file := range files {
		resourceType := inferResourceType(file)
		resourceFiles[resourceType] = append(resourceFiles[resourceType], file)
	}
	return resourceFiles
}
