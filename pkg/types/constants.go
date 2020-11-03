package types

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
)
