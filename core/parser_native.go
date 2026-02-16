//go:build !js && !wasm

package core

import "os"

func ParseFromFilepath(filepath string) (*Farm, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ParseFromReader(file)
}
