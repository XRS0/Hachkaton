package main

import (
	. "hack/services/controller/internal/cli"
	cfg "hack/services/controller/pkg/handleconfig"
)

func main() {
	cfg.ChangeJson(RunCLI())
}
