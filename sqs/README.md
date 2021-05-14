# sqsv2

Sample code to connect to localstack and send a message to an existing SQS queue

We are using the following settings:
```go
queue := "movie-events.fifo"
region := "us-east-1"
endpoint := "http://localhost:4566"
```

## Quick Start

You need to have `localstack-cli` and `aws-cli` installed, optionally also have `jq` installed to format/filter output

1. Run localstack in one termninal
```bash
SERVICES=sqs localstack start
```

2. Create the queue
```bash
aws --endpoint http://localhost:4566 sqs create-queue --queue-name movie-events.fifo
```

3. Run send_msg (sends a message to the existing queue generated in step 2.)
```bash
make run-sqs-send_msg
```

4. Get the message (sent at step 3.) from the queue and cleanup
```bash
# get queue url
aws --endpoint http://localhost:4566 sqs get-queue-url --queue-name movie-events.fifo | jq -r '.QueueUrl'

# receive message from queue (replace the queue url)
aws --endpoint http://localhost:4566 sqs receive-message --queue-url http://localhost:4566/000000000000/movie-events.fifo | jq

# purge queue (replace the queue url)
aws --endpoint http://localhost:4566 sqs purge-queue --queue-url http://localhost:4566/000000000000/movie-events.fifo

# delete queue (replace the queue url)
aws --endpoint http://localhost:4566 sqs delete-queue --queue-url http://localhost:4566/000000000000/movie-events.fifo
```