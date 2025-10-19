package logger

import (
	"fmt"
	"log"
)

func convertInterfaceToString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case int:
		return fmt.Sprintf("%d", v)
	case float64:
		return fmt.Sprintf("%f", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func Debug(param interface{}) {
	log.Println("[DEBUG] " + convertInterfaceToString(param))
}

func Log(param interface{}) {
	log.Println("[LOG] " + convertInterfaceToString(param))
}

func Error(param interface{}) {
	log.Println("[ERROR] " + convertInterfaceToString(param))
}
