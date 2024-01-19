#!/bin/bash
#========================================================================================================#
#         Copied and adapter from `https://github.com/xenon-middleware/xenon-docker-images` :)           #
#========================================================================================================#
echo -e "Starting SLURM..."
# echo -e "starting munged..."
setuser munge /usr/sbin/munged --foreground >/var/log/munged.out.log 2>/var/log/munged.err.log &

# echo -e "starting slurmctld..."
/usr/local/sbin/slurmctld -D -c -vvvv >/var/log/slurmctld.out.log 2>/var/log/slurmctld.err.log &

# echo -e "sleeping for a few seconds to avoid problems..."
sleep 1

# echo -e "starting compute nodes..."
/usr/local/sbin/slurmd -D -N node-0 >/var/log/slurmd-node-0.out.log 2>/var/log/slurmd-node-0.err.log &
/usr/local/sbin/slurmd -D -N node-1 >/var/log/slurmd-node-1.out.log 2>/var/log/slurmd-node-1.err.log &
/usr/local/sbin/slurmd -D -N node-2 >/var/log/slurmd-node-2.out.log 2>/var/log/slurmd-node-2.err.log &
/usr/local/sbin/slurmd -D -N node-3 >/var/log/slurmd-node-3.out.log 2>/var/log/slurmd-node-3.err.log &
/usr/local/sbin/slurmd -D -N node-4 >/var/log/slurmd-node-4.out.log 2>/var/log/slurmd-node-4.err.log &

# echo -e "sleeping for a few seconds to avoid problems..."
sleep 1

# echo -e "making accounting readable to users..."
/bin/chmod -R og+rX /var/log/slurm

# echo -e "starting sshd..."
/usr/sbin/sshd -De >/var/log/sshd.out.log 2>/var/log/sshd.err.log &

echo -e "Startup complete"
#========================================================================================================#