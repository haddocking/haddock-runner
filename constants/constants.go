package constants

const (
	// How many times to wait for a job to finish before giving up,
	//  total time is WAIT_FOR_SLURM * WAIT_TIMEOUT_COUNTER
	//  e.g. 1 * 7200  = 2 hours -> with 1 second sleep
	//  e.g. 1 * 86400 = 24 hours -> with 1 second sleep
	//  e.g. 5 * 7200  = 10 hours -> with 5 second sleep
	WAIT_FOR_SLURM       = 1
	WAIT_TIMEOUT_COUNTER = 86400
)
