package schemaregistry

import (
	"fmt"
	"sort"
	"strings"

	"github.com/confluentinc/go-printer"
)

const (
	SubjectUsage            = "Subject of the schema."
	OnPremAuthenticationMsg = "--ca-location <ca-file-location> --sr-endpoint <schema-registry-endpoint>"
)

var packageDisplayNameMapping = map[string]string{
	"free": "essentials",
	"paid": "advanced",
}

func printVersions(versions []int32) {
	titleRow := []string{"Version"}
	var entries [][]string
	for _, v := range versions {
		record := &struct{ Version int32 }{v}
		entries = append(entries, printer.ToRow(record, titleRow))
	}
	printer.RenderCollectionTable(entries, titleRow)
}

func convertMapToString(m map[string]string) string {
	pairs := make([]string, 0, len(m))
	for key, value := range m {
		pairs = append(pairs, fmt.Sprintf("%s=\"%s\"", key, value))
	}
	sort.Strings(pairs)
	return strings.Join(pairs, "\n")
}

func getPackageDisplayName(packageName string) string {
	return packageDisplayNameMapping[packageName]
}

func getPackageInternalName(packageDisplayName string) (packageInternalName string, isValid bool) {
	packageDisplayName = strings.ToLower(packageDisplayName)
	for k, v := range packageDisplayNameMapping {
		if v == packageDisplayName {
			packageInternalName = k
			isValid = true
			return
		}
	}
	return
}

func getAllPackageDisplayNames() []string {
	packageDisplayNames := make([]string, 0, len(packageDisplayNameMapping))
	for _, tx := range packageDisplayNameMapping {
		packageDisplayNames = append(packageDisplayNames, tx)
	}

	return packageDisplayNames
}
