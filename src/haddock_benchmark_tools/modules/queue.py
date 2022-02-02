import time
import logging

queuelog = logging.getLogger("setuplog")


class Queue:
    def __init__(self, job_list, concurrent=5):
        self.job_list = job_list
        self.concurrent = concurrent
        self.running = []
        # make sure the job_list is sorted
        self.job_list.sort()

    def execute(self):
        """Execute all the jobs in the queue."""
        total = len(self.job_list)
        queuelog.info(
            f"Executing jobs in the queue n={total}, max_concurrent={self.concurrent}"
        )
        while not self.is_done():
            for i, job in enumerate(self.job_list, start=1):
                status = job.status()
                if status == "null" and self.has_slots():
                    queuelog.info(f"> Submitting {job.name} [{i}/{total}]")
                    self.submit(job)

                elif status in ["complete", "failed"]:
                    queuelog.info(f"> Job {job.name} - {status} [{i}/{total}]")
                    self.remove(job)
            time.sleep(60)

    def submit(self, job):
        """Submit for execution."""
        job.run()
        self.running.append(job)

    def remove(self, job):
        """Remove from the running queue."""
        if job in self.running:
            self.running.remove(job)

    def has_slots(self):
        """Check if there are free slots."""
        if len(self.running) < self.concurrent:
            return True
        else:
            return False

    def is_done(self):
        """Check if all the jobs have been executed."""
        if all([e.status() in ["complete", "failed"] for e in self.job_list]):
            return True
        else:
            return False
