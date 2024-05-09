# Use the official Amazon Linux 2 image as the base image
FROM amazonlinux:2

# Update the system packages
RUN yum update -y

# Install git
RUN yum install git -y

# Download and install Go
RUN wget https://go.dev/dl/go1.22.2.linux-amd64.tar.gz && \
    tar -xvf go1.22.2.linux-amd64.tar.gz && \
    mv go /usr/local && \
    rm go1.22.2.linux-amd64.tar.gz

# Set environment variables
ENV GOROOT=/usr/local/go
ENV GOPATH=/app
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH

# Clone the repository
RUN git clone https://github.com/michaelwongycn-work/job-portal.git /app

# Copy the application configuration file
COPY application_config_example.json /app/application_config.json

# Download the Go module dependencies
RUN cd /app && go mod download

# Build the application
RUN cd /app && go build -o app .

# Run the application
CMD ["/app"]