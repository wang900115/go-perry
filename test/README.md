# test

Unit-testing patterns in Go with [`stretchr/testify`](https://github.com/stretchr/testify): assertions, table tests, subtests, mocks, and test suites. The topic has six small packages, each with one function under test.

## What the code does

| Package | Code under test | Test shows |
|---|---|---|
| `func_01` | `Add(a, b)` | a basic `assert.Equal` |
| `func_02` | `IsEven(n)` | a table-driven test over 5 inputs (including 0 and negatives) |
| `func_03` | `Divide(a, b)`, which returns an error when `b == 0` | `require.NoError` / `require.Error` (stops the test on failure) vs `assert` |
| `func_04` | `Greet(name)`, which returns "Hello, Guest" for an empty name | subtests with `t.Run("with name" / "empty name")` |
| `func_05` | `GetUserName(repo UserRepo, id)` | a `mock.Mock` implementation of the `UserRepo` interface, `On(...).Return(...)`, and `AssertExpectations` |
| `func_06` | `GetUsername(repo, id)` (same logic) | a `suite.Suite` with `SetupTest` that creates a fresh mock per test, plus success and "not found" cases |

All tests are black-box tests: they use `package xxx_test` and import the package under test (for example, `test/func_01`).

## Concepts you can learn

- **`assert` vs `require`**: `assert` records a failure and continues. `require` stops the test immediately. Use `require` when the later checks depend on the earlier ones.
- **Table-driven tests**: you add a case by adding a row, not by writing a new test function.
- **Subtests**: named subtests report separately and can be run on their own with `-run TestGreet/empty`.
- **Mocking through interfaces**: because the code depends on `UserRepo`, a test can replace the real data source with a mock. The mock also checks the expected calls.
- **Test suites**: suites group related tests and share setup and teardown (`SetupTest` runs before each test method).
- **External test packages (`_test` suffix)**: these test only the exported API, the way a real caller would use it.

## How to run

```bash
cd test
go test ./...
go test -v -run TestUserTestSuite ./func_06
go test -run 'TestGreet/empty' ./func_04
```

## Notes

- Package names aren't consistent: `func_01` and `func_02` use underscores, while `func03` to `func06` don't, even though the directories are all `func_0N`. Hence the import aliases in the tests.
- `func_05` and `func_06` contain the same production code. The difference is only in how they're tested.
