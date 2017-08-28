# Unseal

Unseal is a small, simple go binary that takes a yaml config file and unseals vault servers.

# Warning

This is currently a WIP, and is not considering production ready or in any way safe.

Use at your own risk

# Why?

When initially deploying vault across multiple sites, you're probably deploying it in a HA config (ie with multiple vault servers in the cluster) and you'll need several people to unseal all of them to get started. This got quite annoying over multiple vault servers and multiple sites, so in order to speed it up, I wrote this little tool.

# Features

Some of the advantages you might gain over using the vault HTTP API or the standard vault binary

  - Zero touch interaction. Once you've written your yaml config, you can simply invoke the command and it'll unseal all the servers that need to be unsealed
  - Parallel execution. Each unseal command runs in a goroutine, meaning you can unseal multiple servers in a matter of seconds
  - Overwriting of unseal key stored in memory. The unseal key you use is zeroed out when the unseal operation is completed, meaning it can't be hijacked by malware etc (see considerations for more info)

# Usage

In order to use unseal, simply create a config file. Here's an example:


```yaml
hosts:
  - name: vault-server-1
    port: 8200
    key: <base64 encoded key>
  - name: vault-server-2
    port: 8200
    key: <base64 encoded key>
  - name: different-site-vault-server.example.com 
    port: 8200
    key: <different base64 encoded key>
```

The app will look for the config file in the following directories, in order:

 - `/etc/unseal/config.yaml`
 - `$HOME/.unseal/config.yaml`
 - `config.yaml` (in the directory you're running the binary from)

Once that's done, simply run the binary:

```bash
./unseal
INFO[0007] Unseal operation performed                    host=site1-consulserver-1 progress=2 threshold=3
INFO[0007] Unseal operation performed                    host=site1-consulserver-2 progress=2 threshold=3
INFO[0008] Unseal operation performed                    host=site1-consulserver-3 progress=2 threshold=3
INFO[0008] Vault is unsealed!                            host=site2-consulserver-2 progress=0 threshold=3
INFO[0008] Vault is unsealed!                            host=site2-consulserver-1 progress=0 threshold=3
INFO[0008] Vault is unsealed!                            host=site2-consulserver-3 progress=0 threshold=3
INFO[0008] Vault is unsealed!                            host=site3-consulserver-1 progress=0 threshold=3
INFO[0008] Vault is unsealed!                            host=site3-consulserver-3 progress=0 threshold=3
INFO[0008] Vault is unsealed!                            host=site3-consulserver-2 progress=0 threshold=3
```

Your vault server progress is now 1 of 3. Yay!

## GPG Support

While you _can_ of course store the unseal keys in plaintext in your `config.yaml` - *it is a really bad idea*. 

With that in mind, Unseal supports GPG decryption. If you've initialized your Vault servers using [PGP/GPG](https://www.vaultproject.io/docs/concepts/pgp-gpg-keybase.html) (and in my opinion, you really _should_) you can specify the base64 encrypted unseal token for your host, and `unseal` will prompt you for your GPG passphrase to decrypt the key.

An example config would look like this:
```
gpg: true
hosts:
  - name: test
  - port: 8200
  - key: <base 64 encoded gpg encrypted key>
```

**Note** - if you have a GPG agent running and you've put the unseal keys in your `config.yaml` - anyone with access to your machine can easily decrypt the values without having to know your GPG password. Be warned.

### Troubleshooting

Unseal simply executes the gpg command to decrypt keys. If you're having any issues with GPG support, I'd suggest doing the following:

1) Ensure you can decrypt the keys manually. Use `echo <base64_key> | base64 -D | gpg -dq`. If this doesn't work, unseal won't work either
2) Ensure you have gpg-agent running, and have a valid `gpg-agent.conf`
3) Ensure your key is a valid base64 encoded string. Again, `echo <base64_key> | base64 -D | gpg -dq` will verify this

## CAPath

Unseal does not support unsecured HTTP API calls, and you probably shouldn't be using Vault over HTTP anyway :)

All your vault servers may use different CA certs, so you can specify a directory with CA certs in it which vault will read and use to attempt to verify the vault server.

Simple specify it like this in your config file:

```yaml
capath: "/path/to/ca/certs"
hosts:
  - name: test
  - port: 8200
  - key: <key>
```

## Environment Variables

By default, vault will read some environment variables to do the unseal config. You can find them [here](https://www.vaultproject.io/docs/commands/environment.html)

You can use _some_ of these environment variables if you wish when using unseal.

 - `VAULT_CACERT`: Set this to the path of a CA Cert you wish to use to verify the vault connection. Note, this will use the same CA cert for all vaults
 - `VAULT_CAPATH`: An alternative to the above CA Path config option.
 - `VAULT_CLIENT_CERT`: An SSL client cert to use when connecting to your vaults. Note, this will use the same cert for all vaults
 - `VAULT_CLIENT_KEY`: An SSL client key to use when connecting to your vaults. Note, this will use the same key for all vaults
 - `VAULT_SKIP_VERIFY`: Skip SSL verification. This is not recommended in production use.

# Considerations

A few security considerations before you use this tool.

 - Your unseal key is clearly stored in plaintext in the yaml file. This is clearly a security issue. Please don't store your unseal key in plaintext permanantly.
 - While I've taken steps to overwrite the unseal key in memory, I am not a Golang expert and it may not be fool proof. If you think you can improve the implementation, pull requests will be warmly welcomed
 - I am just getting started with Golang, and therefore there may be errors, security issues and gremlins in this code. Again, pull requests are much appreciated.
 - There is currently no way of setting HTTPS certificates, so you must trust the certificate presented by the vault API


# Building

If you want to contribute, we use [glide](https://glide.sh/) for dependency management, so it should be as simple as:

 - cloning this repo into `$GOPATH/src/github.com/jaxxstorm/unseal`
 - run `glide install` from the directory
 - run `go build -o unseal main.go`

