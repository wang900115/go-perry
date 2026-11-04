# dump

Shows that secrets held in a process's memory can be read from a memory dump of that process.

## What the code does

`main.go` stores `secret := "prod-db-password-123"`, prints `service started` and then the secret, and sleeps for one hour so there is time to dump the running process.

The directory also contains committed artifacts:

- `dump.exe`: the built Windows binary.
- `dump.exe_260726_145304.dmp` (~27 MB): a Windows minidump (the file starts with `MDMP`). The literal string `prod-db-password-123` appears in it twice. The naming pattern suggests a tool such as ProcDump created it, but nothing in the repo confirms this.

## Concepts you can learn

- **Secrets live in plaintext in memory**: anyone who can dump the process (an admin, a crash reporter, a core dump) can read string constants and loaded credentials with a simple text search.
- **String literals live in the binary**: `prod-db-password-123` is compiled into the executable itself, not only into the heap at runtime.
- **Long-running processes are dump targets**: the one-hour sleep stands in for a service that keeps secrets in memory for its whole lifetime.

## How to run

```bash
cd dump
go run .     # blocks for 1 hour; Ctrl+C to stop
```

To reproduce the dump, create a dump of the running process with an OS tool (for example Task Manager "Create dump file" or ProcDump on Windows, or `gcore` on Linux). Then search it with `strings` or `grep -a`.

## Notes

- The program does not try to protect the secret. It only demonstrates the exposure.
- The committed `.dmp` file is large and contains the demo secret. Do not commit dumps of real services.
