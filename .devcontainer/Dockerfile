#========================================================================================================#
FROM xenonmiddleware/slurm:20

WORKDIR /app
COPY start.sh /app/start.sh
RUN chmod +x /app/start.sh

COPY --from=golang:1.21 /usr/local/go/ /usr/local/go/
#========================================================================================================#
