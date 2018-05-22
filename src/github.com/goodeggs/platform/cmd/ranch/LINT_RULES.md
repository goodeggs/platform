
## cron-minutes

On Ranch, cron jobs should run at intervals >= 10 minutes, or at a specific minute.  If you need a job that runs with greater frequency, you should write it as a long-lived worker process (see https://github.com/goodeggs/goodeggs-server#script-runner for guidance).  We require intervals >= 10 minutes because the lifecycle of each cron job takes around 10 minutes, and therefore has a risk of job overrun and resource exhaustion.

