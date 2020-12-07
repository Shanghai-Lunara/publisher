package types

import (
	"fmt"
	"time"
)

const (
	// ws
	WebsocketHandlerRunner    = "/runner"
	WebsocketHandlerDashboard = "/dashboard"

	PublisherProjectDir = "PUBLISHER_PROJECT_DIR"
	// git config
	PublisherGitBranch = "PUBLISHER_GIT_BRANCH"

	// ftp config
	PublisherFtpHost     = "ftp_host"
	PublisherFtpPort     = "ftp_port"
	PublisherFtpUsername = "ftp_username"
	PublisherFtpPassword = "ftp_password"
	PublisherFtpWorkDir  = "ftp_work_dir"
	PublisherFtpTimeout  = "ftp_timeout"
	PublisherFtpMkdir    = "PUBLISHER_FTP_MKDIR"

	// svn config
	PublisherSvnHost          = "svn_host"
	PublisherSvnPort          = "svn_port"
	PublisherSvnUsername      = "svn_username"
	PublisherSvnPassword      = "svn_password"
	PublisherSvnRemoteDir     = "svn_remote_dir"
	PublisherSvnWorkDir       = "svn_work_dir"
	PublisherSvnCommitMessage = "svn_commit_message"
	PublisherSvnCommand       = "svn_command"

	// version flag
	VersionFlag = "VersionFlag"

	// Robot
	RobotDurationInMs = "Robot_Duration_In_MS"
)

const (
	RecordDefault = iota
	RecordVersion
)

const (
	StepMessageFormat = "[%s] StepName: [%s] Message: [%s] is starting"
)

func StepMessage(stepName, action string) string {
	return fmt.Sprintf(StepMessageFormat, time.Now().Format("2006-01-02 15:04:05"), stepName, action)
}
