package main

import "os"

// Simple check for whather the file exists.
// TODO: be more robust, i.e. permissions etc.
func FileExists(filename string) (exists bool, err error) {
	_, err = os.Stat(filename)
	if err != nil {
		return false, nil
	}
	return true, nil
}
