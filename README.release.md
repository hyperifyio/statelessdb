## What is it?

StatelessDB is a cloud database server that processes encrypted client-side data 
without storing state, enabling secure and scalable computations.

## Running the server

You can run the server manually:

```bash
./statelessdb
```

Once started, the server is available at http://localhost:8080/api/v1

## Installing StatelessDB

Install the server executable:

```bash
cp ./statelessdb /usr/local/bin/statelessdb
chmod 700 /usr/local/bin/statelessdb
chown nobody:nogroup /usr/local/bin/statelessdb
```

For systemd configuration, install the provided service file to
`/etc/systemd/system/statelessdb.service`.

Use `/usr/local/bin/statelessdb --init-private-key` to generate a private key for the 
server.

Configure the private key and PostgreSQL settings to the env configuration to 
`/etc/statelessdb/env`.
