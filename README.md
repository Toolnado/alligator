# Alligator

**Alligator** is simple distributed in-memory cache written in go.

## Features

- In-memory caching for fast data retrieval
- Supports connection via TCP protocol
- Simple and intuitive API
- Thread-safe operations

## Installation

To install Alligator, use `git clone`:

```bash
git clone github.com/Toolnado/alligator
```

## Usage

Here's a basic example of how to use # Alligator:

1. Сreate a binary using makefile:
```bash
make compile
```
2. Launch the leader:
```bash
./build/app_amd64 --addr :3000
```
3. Launch the followers:
```bash
./build/app_amd64 --addr :4000 --laddr :3000
```

## Configuration

You can configure the cache using flags:

1. --addr - listener address;
2. --laddr - cache leader address.

If the leader address is not specified, then the current instance will be the leader.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes. Be sure to follow the project's coding guidelines and include tests for any new features or bug fixes.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

