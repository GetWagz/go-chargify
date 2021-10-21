package utils

import (
	"encoding/json"

	"github.com/spf13/cobra"
)

func ParentPersistentPreRunE(cmd *cobra.Command, args []string) error {
	parent := cmd.Parent()
	if parent != nil {
		if parent.PersistentPreRunE != nil {
			err := parent.PersistentPreRunE(parent, args)
			if err != nil {
				return err
			}
		} else {
			ParentPersistentPreRunE(parent, args)
		}
	}

	return nil
}

//PrettyJSON from object
func PrettyJSON(obj interface{}) string {
	jsonBytes, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}

// JSON from object
func JSON(obj interface{}) string {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}
