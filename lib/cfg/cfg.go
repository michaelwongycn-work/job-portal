package cfg

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"

	"github.com/michaelwongycn/job-portal/domain/config"
	"github.com/michaelwongycn/job-portal/lib/json"
)

const fname = "application_config.json"

func ReadConfig() (*config.ApplicationConfig, error) {
	flag.Parse()

	var cfg config.ApplicationConfig

	current_dir := "."
	pathsep := string(os.PathSeparator)
	for i := 0; i < 255; i++ {
		err := json.ReadJSON(current_dir+pathsep+fname, &cfg)
		if err == nil {
			return &cfg, nil
		}

		if !errors.Is(err, fs.ErrNotExist) {
			return nil, fmt.Errorf("error parsing file: %v", err)
		}
		current_dir = ".." + pathsep + current_dir
	}
	return nil, fmt.Errorf("config file %s not found", fname)
}
