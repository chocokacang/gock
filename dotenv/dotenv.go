package dotenv

import (
	"bytes"
	"flag"
	"io"
	"os"

	"github.com/chocokacang/gock/log"
)

type envVars map[string]string

func Load() {
	env := flag.String("env", "LOCAL", "Set application environment: PRODUCTION, LOCAL, TEST")
	flag.Parse()

	os.Setenv("APP_ENV", *env)

	envMap := read()
	for key, value := range envMap {
		os.Setenv(key, value)
	}
}

func filename() string {

	switch os.Getenv("APP_ENV") {
	case "LOCAL":
		return ".env.local"
	case "TEST":
		return ".env.test"
	default:
		return ".env"
	}
}

func read() envVars {
	file, err := os.Open(filename())
	if err != nil {
		log.Warning("Failed to %v", err)
		return nil
	}
	defer file.Close()

	return Parse(file)
}

func Parse(r io.Reader) envVars {
	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		log.Warning("Failed to parse env variable. Got error: %v", err)
	}

	parsed, _ := UnmarshalBytes(buf.Bytes())

	return parsed
}

// UnmarshalBytes parses env file from byte slice of chars, returning a map of keys and values.
func UnmarshalBytes(src []byte) (map[string]string, error) {
	out := make(map[string]string)
	err := parseBytes(src, out)

	return out, err
}
