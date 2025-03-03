package fidoctl_test

import (
	"fmt"
	"log"

	"github.com/buglloc/fidoctl"
)

func ExampleEnumerate() {
	devices, err := fidoctl.Enumerate()
	if err != nil {
		panic(fmt.Errorf("Enumerate: %v", err))
	}

	for _, device := range devices {
		fmt.Println(device.String())
		cfg, err := device.YubiConfig()
		if err != nil {
			log.Printf("Getting YubiConfig: %v", err)
			continue
		}

		fmt.Printf("  - serial: %d\n", cfg.Serial())
		fmt.Printf("  - version: %s\n", cfg.Version())
	}
}
