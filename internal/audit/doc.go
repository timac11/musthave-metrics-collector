/*
This package provides functionality to write audit logs to files and server

To create auditor:
	auditor := audit.NewAuditor(auditsPath, auditsURL)

	where:
		auditsPath - path to json file, where logs will be written
		auditsURL - url where audit logs will be sended
*/

package audit
