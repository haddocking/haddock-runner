```
sudo pacman -S slurm-llnl
sudo systemctl enable --now munge
sudo mkdir -p /etc/slurm-llnl
sudo cp cgroup.conf /etc/slurm-llnl/
sudo cp slurm.conf /etc/slurm-llnl/
sudo systemctl enable --now slurmctld
sudo systemctl enable --now slurmd
```

```
sinfo
squeue
srun --nodes=1 --ntasks=1 hostname
```


