# pub-dev

A self-hosted Dart and Flutter package registry, written in Go. This server implements the [Dart Pub Repository Specification](https://github.com/dart-lang/pub/blob/master/doc/repository-spec-v2.md) and allows you to host your own private or public packages.

## Features

-   **Self-Hosted:** Host your own packages on your own infrastructure.
-   **File-Based Storage:** Packages are stored directly on the filesystem, making it simple to manage and back up.
-   **Authentication:** Secure your server with token-based authentication for package publishing.
-   **CLI Tools:** Includes a command-line interface for server administration.

## Getting Started

### Prerequisites

-   [Go](https://golang.org/dl/) (version 1.21 or later)

### Building from source

1.  Clone the repository:
    ```sh
    git clone https://github.com/fmotalleb/pub-dev.git
    cd pub-dev
    ```

2.  Build the binary:
    ```sh
    go build .
    ```

## Usage

### Running the Server

To run the server, you need a configuration file. An example is provided in `example/config.toml`.

```sh
./pub-dev --config /path/to/your/config.toml
```

### Command-Line Interface

The application is controlled via the `pub-dev` command.

#### Global Flags

-   `-c`, `--config`: Path to the configuration file (default: `./config.toml`).
-   `-v`, `--verbose`: Enable verbose output for debugging.

#### Commands

-   **`run`** (default): Starts the pub server.
    ```sh
    ./pub-dev
    ```

-   **`re-calculate`**: This command regenerates the `listing.json` metadata file for all packages within the storage directory. This is useful if the metadata files become corrupted or need to be recreated.
    ```sh
    ./pub-dev re-calculate -s /path/to/storage
    ```

## Configuration

Configuration is managed through a TOML file.

**Example `config.toml`:**

```toml
# The address and port the HTTP server will listen on.
http_listen = ":8080"

# The base URL of the server, used for generating package URLs.
# This should be the public-facing URL.
base_url = "http://localhost:8080/"

# The path to the directory where packages will be stored.
storage = "./storage/packages"

# Optional: Define authentication rules for specific API paths.
[[auth]]
  # A list of URL paths to protect.
  path = ["/api/packages/versions/newUpload"]
  # A list of valid bearer tokens for these paths.
  token = ["your-secret-token-here"]
```

### Configuration Options

-   `http_listen` (String): The TCP address for the server to listen on.
-   `base_url` (String): The public URL for the server. This is critical as it's used to construct the `archive_url` for packages.
-   `storage` (String): The local filesystem path where packages will be stored.
-   `auth` (Array of Tables): Defines authentication rules.
    -   `path` (Array of Strings): A list of URL prefixes to protect.
    -   `token` (Array of Strings): A list of allowed bearer tokens for the specified paths.

## Publishing Packages

To publish a package to your `pub-dev` server, you need to override the default pub server in your package's `pubspec.yaml`:

```yaml
name: my_awesome_package
description: An awesome package.
version: 1.0.0

environment:
  sdk: '>=3.0.0 <4.0.0'

# Add this section to publish to your private server
publish_to: http://your-server-address:8080

dependencies:
  flutter:
    sdk: flutter
```

Then, you can publish it using the standard Dart or Flutter CLI. If you have authentication configured, you will need to set up authentication on your local machine.

```sh
# For Dart packages
dart pub publish

# For Flutter packages
flutter pub publish
```

## API Endpoints

The server exposes the following main API endpoints:

-   `GET /api/packages/:package`: Retrieves the metadata for a specific package.
-   `GET /api/packages/versions/new`: Initiates the package upload process.
-   `POST /api/packages/versions/newUpload`: Uploads the package `.tar.gz` archive.
-   `GET /api/packages/versions/newUploadFinish`: Finalizes the upload.
-   `GET /storage/packages/...`: Serves the package `.tar.gz` files.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.
