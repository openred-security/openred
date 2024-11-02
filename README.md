# openred

openred is a fully open-source security platform that aggregates a curated catalog of open-source security tools, making them easily accessible and manageable. With openred, users can leverage a broad range of external security tools through modular plugins for tasks such as vulnerability scanning, system inventory collection, compliance reporting, and file signature analysis. Each plugin is a separate open-source project that openred has reviewed, adapted, and standardized for seamless integration.

## Key Features

- **Plugin Catalog**: A growing list of curated open-source security tools, each adapted for straightforward usage within openred.
- **Unified Interface**: Standardized input and output for all plugins, offering consistency regardless of the tool.
- **CLI for Easy Execution**: Run plugins with simple commands, no need for in-depth knowledge of each tool's intricacies.
- **Flexible Security Solutions**: Empower users to tailor their security setup by choosing from a wide range of security plugins in the catalog.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/openred.git
   cd openred
   ```

2. Build the project using Go:
   ```bash
   go build -o openred main.go
   ```

3. Run openred:
   ```bash
   ./openred run <plugin_name>
   ```

## Usage

Execute openred with any plugin available in the catalog. For example:

```bash
./openred run <plugin_name>
```

### Example: Running the `last_logs` Plugin

To run the `last_logs` plugin:

```bash
./openred run last_logs
```

## Plugin Catalog

The catalog consists of external open-source projects that openred integrates, adapts, and normalizes to make them usable with minimal configuration. Each plugin is configured through a `config.yml` file in the plugin's directory, where command execution and output handling are defined.

### Plugin Example: `last_logs`

- `config.yml`:
   ```yaml
   plugin_name: last_logs
   executable: "last_logs.sh"
   ```

## TODO / What's Next

- **openred Console Integration**: We’re developing an interface, openred Console, to centralize and visualize plugin outputs. The console is based on OpenSearch and OpenSearch Dashboards but with a simplified UI to enhance usability.
- **Future SaaS Version**: openred aims to offer a SaaS version to simplify security management for organizations of all sizes.

## Contributing

Contributions are welcome! Please refer to the [Contributing Guide](CONTRIBUTING.md) for guidelines on how to get started.

## License

openred is licensed under the [GNU General Public License v3.0](LICENSE).
