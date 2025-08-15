# Protun

**protun = HTTP Proxy + OpenVPN tunnel**

A Go-based application that establishes an OpenVPN connection and provides an HTTP proxy server that routes requests through the VPN tunnel. Perfect for routing HTTP traffic through specific VPN endpoints with containerized deployment support.

## Overview

**protun** combines two main components:
- **VPN Client**: Establishes and maintains OpenVPN connections using `.ovpn` configuration files
- **HTTP Proxy Server**: Gin-based HTTP server that proxies requests through the established VPN tunnel

The application automatically handles DNS configuration updates and provides logging for both VPN connection status and proxy requests.

## Features

- 🔒 **Secure VPN Integration**: Uses OpenVPN3 library for reliable VPN connections
- 🚀 **HTTP Proxy Server**: Built with Gin framework for high-performance request proxying  
- 🐳 **Docker Support**: Complete containerization with Docker and Docker Compose
- 📊 **Comprehensive Logging**: Detailed logging for VPN events, connection status, and proxy requests
- ⚙️ **Flexible Configuration**: Environment-based configuration for easy deployment
- 🌍 **Multi-Location Support**: Easy setup for multiple VPN endpoints

## Prerequisites

### Local Development
- **Go 1.23+**
- **OpenVPN** installed on your system
- **Root/Administrator privileges** (required for VPN operations and DNS updates)

### Docker Deployment  
- **Docker** and **Docker Compose**
- **Privileged container access** (required for VPN operations)

## Project Structure

```
protun/
├── main.go                 # Application entry point
├── proxy_server/
│   └── proxy_server.go     # HTTP proxy server implementation
├── vpn/
│   ├── vpn.go             # VPN connection management
│   └── logging.go         # VPN logging and DNS configuration
├── ovpn/
│   └── configs/           # Directory for .ovpn configuration files
├── go.mod                 # Go module dependencies
├── Dockerfile             # Docker container configuration
└── docker-compose.yml     # Multi-service Docker setup
```

## Installation

**Clone the repository**
   ```bash   
   git clone https://github.com/leBruchette/protun.git
   cd protun
   ```
### Method 1: Local Go Installation

**Install dependencies**
   ```bash
   go mod tidy
   ```

**Build the application**
   ```bash
   go build -o protun main.go
   ```

### Method 2: Docker Build

   ```bash
   docker build --platform linux/amd64 -t protun .
   ```

## Configuration

### Environment Variables

The application requires the following environment variables:

| Variable | Description | Example         |
|----------|-------------|-----------------|
| `VPN_CONFIG` | Name of the .ovpn file (without extension) | `us_oregon`     |
| `VPN_USER` | OpenVPN username | `your_username` |
| `VPN_PASS` | OpenVPN password | `your_password` |

### VPN Configuration Files

1. **Create the configs directory**
   ```bash
   mkdir -p ovpn/configs
   ```

2. **Add your .ovpn files**
   - Place your OpenVPN configuration files in `ovpn/configs/`
   - Name them according to your `VPN_CONFIG` environment variable
   - Example: `ovpn/configs/oregon.ovpn`

3. **Ensure proper file permissions**
   ```bash
   chmod 755 ovpn/configs
   chmod 644 ovpn/configs/*.ovpn
   ```

## Usage

### Running Locally

1. **Set environment variables**
   ```bash
   export VPN_CONFIG="us_oregon"
   export VPN_USER="your_username"  
   export VPN_PASS="your_password"
   ```

2. **Run with elevated privileges**
   ```bash
   sudo ./protun
   ```

### Running with Docker

#### Single Container
```bash
docker run --privileged \
  -e VPN_CONFIG="us_oregon" \
  -e VPN_USER="your_username" \
  -e VPN_PASS="your_password" \
  -v /path/to/your/ovpn/configs:/app/ovpn \
  -p 80:8080 \
  protun
```

#### Docker Compose (Recommended)

1. **Update docker-compose.yml** with your credentials:
   ```yaml
   environment:
     - VPN_USER=your_actual_username
     - VPN_PASS=your_actual_password
   volumes:
     - /path/to/your/ovpn/configs:/app/ovpn/configs
   ```
2. **Start all services**
   ```bash
   docker-compose up 
   ```


3. **Start specific service**
   ```bash
   # Start US Oregon proxy on port 80
   docker-compose up us-oregon
   
   # Start US Texas proxy on port 81  
   docker-compose up us-texas
   
   # Start US Minnesota proxy on port 82
   docker-compose up us-minnesota
   ```

## API Usage

Once running, the proxy server accepts HTTP requests on the configured port (default: 8080).

### Proxy Endpoint

**Endpoint**: `ANY /proxy/*path`

**Method**: All HTTP methods (GET, POST, PUT, DELETE, etc.)

**Example Requests**:
```bash
# GET request through proxy
curl http://localhost:8080/proxy/example.com/api/data

# POST request through proxy  
curl -X POST http://localhost:8080/proxy/api.example.com/endpoint \
  -H "Content-Type: application/json" \
  -d '{"key": "value"}'
```

**Note**: Currently, the proxy routes all requests to a fixed endpoint for testing. Update the `fullURL` variable in `proxy_server.go` to customize the target destination.

## Application Flow

1. **VPN Connection**: Application starts VPN session using provided credentials and configuration
2. **DNS Update**: Automatically updates system DNS to use VPN-provided DNS servers  
3. **Proxy Server**: Starts HTTP proxy server on port 8080
4. **Request Handling**: All HTTP requests to `/proxy/*` are forwarded through the VPN tunnel
5. **Logging**: Comprehensive logging of VPN events and proxy requests

## Logging

The application provides detailed logging for:

- **VPN Connection Events**: Connection status, authentication, and tunnel establishment
- **DNS Updates**: Automatic DNS server configuration changes
- **Proxy Requests**: HTTP request details and response forwarding
- **Performance Statistics**: Periodic VPN connection statistics

## Troubleshooting

### Common Issues

**Permission Denied**
- Ensure you're running with `sudo` or administrator privileges
- Docker containers need `--privileged` flag

**VPN Connection Fails**
- Verify your `.ovpn` file is valid and accessible
- Check that `VPN_USER` and `VPN_PASS` are correct
- Ensure the `.ovpn` file path matches your `VPN_CONFIG` environment variable

**DNS Resolution Issues**
- The application automatically updates `/etc/resolv.conf`
- Verify the container/system allows DNS configuration changes

**Port Already in Use**
- Change the port mapping in Docker or the port in `proxy_server.go`
- Check for other services using port 8080

### Debug Mode

Enable verbose logging by modifying the `disableLogs` variable in `logging.go`:

```go
var disableLogs = false // Keep logs enabled for debugging
```

## Security Considerations

- **Privileged Access**: Required for VPN operations and DNS modifications
- **Credential Storage**: Environment variables contain sensitive VPN credentials
- **TLS Verification**: Currently disabled (`InsecureSkipVerify: true`) for proxy requests
- **Network Isolation**: Consider running in isolated network environments

## Development

### Adding New VPN Endpoints

1. Add new `.ovpn` configuration file to `ovpn/configs/`
2. Add new service in `docker-compose.yml`:
   ```yaml
   new-location:
     extends:
       service: common
     environment:
       - VPN_CONFIG=new_location_config
     ports:
       - "83:8080"
   ```

### Customizing Proxy Behavior

Modify `proxy_server.go` to customize:
- Target URL routing logic
- Request/response header handling  
- Authentication mechanisms
- Request logging and filtering

## Dependencies

- **[github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)**: HTTP web framework
- **[github.com/mysteriumnetwork/go-openvpn](https://github.com/mysteriumnetwork/go-openvpn)**: OpenVPN 3 Go bindings
- **Standard Go libraries**: net/http, crypto/tls, os, sync

## License

[Add your license information here]

## Contributing

[Add contribution guidelines here]