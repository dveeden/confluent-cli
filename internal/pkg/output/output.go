package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/confluentinc/go-printer"
	"github.com/go-yaml/yaml"
	"github.com/spf13/cobra"
	"github.com/tidwall/pretty"

	"github.com/confluentinc/cli/internal/pkg/errors"
)

const (
	humanString   = "human"
	jsonString    = "json"
	yamlString    = "yaml"
	FlagName      = "output"
	ShortHandFlag = "o"
	Usage         = `Specify the output format as "human", "json", or "yaml".`
	DefaultValue  = humanString
)

var (
	allFormatStrings = []string{humanString, jsonString, yamlString}
)

type Format int

const (
	Human Format = iota
	JSON
	YAML
)

func (o Format) String() string {
	return allFormatStrings[o]
}

type ListOutputWriter interface {
	/*
		AddElement - Add an element to the list to output for StructuredListWriter
		* @param e : the element to add, must be either a pointer or an interface
	*/
	AddElement(e interface{})
	/*
		Out - Create the output to the IO channel passed in during construction
	*/
	Out() error
	GetOutputFormat() Format
	StableSort()
}

/*
NewListOutputWriter - Create a new ListOutputWriter.
Returns an ListWriter that is used to output a list of objects (must be pointers of an interface) in different formats (humanreadable, json, yaml)
 * @param cmd: The cobra.Command called
 * @param listFields: A list of fields (of the underlying object we're outputting) that we want to output
 * @param humanLabels: A list of names for the fields (n the same order) that we want in the output for the human readable view
 * @param structedLabels: A list of names for the fields (in the same order) that we want in the output for structured views (yaml and json)
@return ListOutputWriter, error
*/
func NewListOutputWriter(cmd *cobra.Command, listFields []string, humanLabels []string, structuredLabels []string) (ListOutputWriter, error) {
	return NewListOutputCustomizableWriter(cmd, listFields, humanLabels, structuredLabels, cmd.OutOrStdout())
}

func NewListOutputCustomizableWriter(cmd *cobra.Command, listFields []string, humanLabels []string, structuredLabels []string, writer io.Writer) (ListOutputWriter, error) {
	if len(listFields) != len(humanLabels) || len(humanLabels) != len(structuredLabels) {
		return nil, errors.New("argument list length mismatch") // TODO: correct error to return?
	}
	format, err := cmd.Flags().GetString(FlagName)
	if err != nil {
		return nil, err
	}
	switch format {
	case JSON.String():
		return &StructuredListWriter{
			outputFormat: JSON,
			listFields:   listFields,
			listLabels:   structuredLabels,
			writer:       writer,
		}, nil
	case YAML.String():
		return &StructuredListWriter{
			outputFormat: YAML,
			listFields:   listFields,
			listLabels:   structuredLabels,
			writer:       writer,
		}, nil
	case Human.String():
		return &HumanListWriter{
			outputFormat: Human,
			listFields:   listFields,
			listLabels:   humanLabels,
			writer:       writer,
		}, nil
	default:
		return nil, NewInvalidOutputFormatFlagError(format)
	}
}

func DescribeObject(cmd *cobra.Command, obj interface{}, fields []string, humanRenames, structuredRenames map[string]string) error {
	format, err := cmd.Flags().GetString(FlagName)
	if err != nil {
		return err
	}
	if !(format == Human.String() || format == JSON.String() || format == YAML.String()) {
		return NewInvalidOutputFormatFlagError(format)
	}
	return printer.RenderOut(obj, fields, humanRenames, structuredRenames, format, os.Stdout)
}

// StructuredOutput - pretty prints an object in specified format (JSON or YAML) using tags specified in struct definition
func StructuredOutput(format string, obj interface{}) error {
	var b []byte
	if format == JSON.String() {
		j, _ := json.Marshal(obj)
		b = pretty.Pretty(j)
	} else if format == YAML.String() {
		b, _ = yaml.Marshal(obj)
	} else {
		return NewInvalidOutputFormatFlagError(format)
	}
	_, err := fmt.Fprintf(os.Stdout, string(b))
	return err
}

// NewInvalidOutputFormatFlagError - create a new error to describe an invalid output format flag
func NewInvalidOutputFormatFlagError(format string) error {
	errorMsg := fmt.Sprintf(errors.InvalidFlagValueErrorMsg, format, FlagName)
	suggestionsMsg := fmt.Sprintf(errors.InvalidFlagValueSuggestions, FlagName, strings.Join(allFormatStrings, ", "))
	return errors.NewErrorWithSuggestions(errorMsg, suggestionsMsg)
}

// IsValidFormatString - returns whether a format string is a valid format (human, json, yaml)
func IsValidFormatString(format string) bool {
	for _, formatString := range allFormatStrings {
		if format == formatString {
			return true
		}
	}
	return false
}
