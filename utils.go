package lambroll

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Songmu/prompter"
	"github.com/aws/aws-sdk-go/private/protocol/json/jsonutil"
)

func (app *App) saveFile(path string, b []byte, mode os.FileMode) error {
	if _, err := os.Stat(path); err == nil {
		ok := prompter.YN(fmt.Sprintf("Overwrite existing file %s?", path), false)
		if !ok {
			return nil
		}
	}
	return ioutil.WriteFile(path, b, mode)
}

func marshalJSON(s interface{}) ([]byte, error) {
	var buf bytes.Buffer
	b, err := jsonutil.BuildJSON(s)
	if err != nil {
		return nil, err
	}
	json.Indent(&buf, b, "", "  ")
	buf.WriteString("\n")
	return buf.Bytes(), nil
}

func marshalJSONV2(s interface{}) ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

func unmarshalJSON(src []byte, v interface{}, path string) error {
	strict := json.NewDecoder(bytes.NewReader(src))
	strict.DisallowUnknownFields()
	if err := strict.Decode(&v); err != nil {
		if !strings.Contains(err.Error(), "unknown field") {
			return err
		}
		log.Printf("[warn] %s in %s", err, path)

		// unknown field -> try lax decoder
		lax := json.NewDecoder(bytes.NewReader(src))
		return lax.Decode(&v)
	}
	return nil
}

func FindFunctionFilename() string {
	for _, name := range FunctionFilenames {
		if _, err := os.Stat(name); err == nil {
			return name
		}
	}
	return FunctionFilenames[0]
}
