use crate::job::Job;
use anyhow::Context;
use log::{error, info};
use std::sync::Arc;
use std::sync::mpsc;
use std::thread;

pub struct Queue {
    concurrent: u16,
    workload: Vec<Job>,
}

impl Queue {
    pub fn new(concurrent: u16, workload: Vec<Job>) -> Self {
        Queue {
            concurrent,
            workload,
        }
    }

    pub fn setup(&self) -> anyhow::Result<()> {
        // Run setup for all jobs sequentially - this is I/O and doesn't benefit much from parallelism
        for job in &self.workload {
            let mut job_clone = job.clone();
            info!("Setting up job: {}", job_clone.name);
            job_clone.setup()?;
        }
        Ok(())
    }

    pub fn start(&self) -> anyhow::Result<()> {
        // Create a thread pool with the specified concurrency level
        let (tx, rx) = mpsc::channel();
        let workload = Arc::new(self.workload.clone());
        let concurrent = self.concurrent as usize;

        // Spawn worker threads
        for i in 0..concurrent {
            let tx = tx.clone();
            let workload = workload.clone();
            let thread_id = i;

            thread::spawn(move || {
                // Process jobs in this thread
                for (index, job) in workload.iter().enumerate() {
                    if index % concurrent == thread_id {
                        let mut job_clone = job.clone();
                        info!("Thread {} processing job: {}", thread_id, job_clone.name);

                        if let Err(e) = job_clone.run() {
                            error!(
                                "Thread {} failed processing job {}: {}",
                                thread_id, job_clone.name, e
                            );
                            tx.send(Err(e)).unwrap();
                            return;
                        }
                    }
                }

                tx.send(Ok(())).unwrap();
            });
        }

        // Wait for all threads to complete and collect results
        let mut results = Vec::new();
        for _ in 0..concurrent {
            results.push(
                rx.recv()
                    .context("Failed to receive result from worker thread")?,
            );
        }

        // Check if any thread failed
        for result in results {
            result?;
        }

        Ok(())
    }
}
