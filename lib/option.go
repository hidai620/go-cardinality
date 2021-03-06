package dbindex

import (
	"errors"
	"flag"
	"fmt"
	sutil "github.com/hidaiy/go-utils/stringutil"
	"os"
	"strings"
)

// OutputType
type OutputType int

const (
	_ OutputType = iota
	CONSOLE
	CSV
)

// GetOutPutType is factory function of OutputType.
// It returns OutputType from string.
func GetOutputType(o string) (OutputType, error) {
	switch strings.ToUpper(o) {
	case "CONSOLE":
		return CONSOLE, nil
	case "CSV":
		return CSV, nil
	default:
		return 0, errors.New(fmt.Sprintf("Pselese select type following list. console, csv"))
	}
}

// String returns OutputType as human readable string.
func (o OutputType) String() string {
	var ret string
	switch o {
	case CONSOLE:
		ret = "CONSOLE"
	case CSV:
		ret = "CSV"
	}
	return ret
}

// Parse parses command line flags.
func ParseCommandLineOption() (*Option, error) {
	var (
		config, out, tableNames string
		allTable                bool
	)
	flag.StringVar(&config, "config", "config.toml", "Absolute or relrative path of config file.")
	flag.StringVar(&out, "out", "console", `Output type of result. "console" or "csv"`)
	flag.StringVar(&tableNames, "table", "", `Analyze Target table name.`)
	flag.BoolVar(&allTable, "allTable", false, `Analyze all table.`)
	flag.Parse()

	// check where config file exists.
	err := existsFile(config)
	if err != nil {
		return nil, err
	}

	// table flag check
	if isBothTableOptionSpecified(allTable, tableNames) {
		return nil, errors.New("You can specity table or allTable.")
	}

	if isNotSpecifiedTables(allTable, tableNames) {
		return nil, errors.New("Prease specity table or allTable.")
	}

	// get OutputType
	outputType, err := GetOutputType(out)
	if err != nil {
		return nil, err
	}

	tableNamesArray := sutil.Split(tableNames, ",")
	return NewOption(outputType, config, tableNamesArray), nil
}

// isBothTableOptionSpecified returns true if allTable is true, and tableNames is not empty.
func isBothTableOptionSpecified(allTable bool, tableNames string) bool {
	return allTable && tableNames != ""
}

// isNotSpecifiedTables returns true argument does not have some values.
func isNotSpecifiedTables(allTable bool, tableNames string) bool {
	return !allTable && tableNames == ""
}

// Option has command line arguments.
type Option struct {
	Out        OutputType
	ConfigPath string
	TableNames []string
}

// New returns CommandLineOption created with arguments
func NewOption(out OutputType, configPath string, tableNames []string) *Option {
	return &Option{
		Out:        out,
		ConfigPath: configPath,
		TableNames: tableNames,
	}
}

// Equals returns true if argument is save structure, and has same value.
func (c *Option) Equals(c2 *Option) bool {
	return c.Out == c2.Out &&
		c.ConfigPath == c2.ConfigPath
}

// existsFile returns error if argument filePath does not exist.
func existsFile(filePath string) error {
	_, err := os.Stat(string(filePath))
	if err != nil {
		return errors.New(fmt.Sprintf("Specified file does not exist. %s", filePath))
	}
	return nil
}
