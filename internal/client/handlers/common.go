package handlers

import (
	"fmt"

	"github.com/githubzjm/tuo/api/v1/def"
	"github.com/githubzjm/tuo/internal/client/common"
)

// Welcome to the MySQL monitor.  Commands end with ; or \g.
// Your MySQL connection id is 19
// Server version: 8.0.28 MySQL Community Server - GPL

// Copyright (c) 2000, 2022, Oracle and/or its affiliates.

// Oracle is a registered trademark of Oracle Corporation and/or its
// affiliates. Other names may be trademarks of their respective
// owners.

// Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.

func handleRespErr(status int, resp *def.BaseResp, err error) bool {
	if err != nil {
		fmt.Printf(common.ErrorColor("INTERNAL ERROR: %s\n"), err.Error())
		// common.GracefulExit(1)
		return false
	} else if status > 299 || status < 0 {
		fmt.Printf(common.ErrorColor("STATUS: %d\nERROR: %s\n"), status, resp.Error)
		return false
	}
	return true
}
