# Contributing to go-perry

Thank you for contributing to `go-perry`.

`go-perry` is a collection of Go experiments and topic-based implementations.

Contributions are not limited to fixing bugs or improving existing examples. You are encouraged to create your own concept and turn it into a new topic.

## Contribution Philosophy

The main idea is simple:

> **Create your own concept. Implement it. Document it. Share it.**

You can contribute:

- A new Go topic
- A technical experiment
- A new implementation
- A benchmark
- A runtime experiment
- A concurrency example
- A networking example
- A performance experiment
- A systems programming experiment
- Improvements to existing topics

If you have an idea worth exploring, create a feature branch and experiment with it.

---

## 1. Create a Feature Branch

Do not work directly on the `main` branch.

Fork the repository, then create a feature branch:

```bash
git checkout main
git pull origin main
git checkout -b feature/<your-concept>
```

For example:

```bash
git checkout -b feature/tracing
```

```bash
git checkout -b feature/goroutine-scheduler
```

```bash
git checkout -b feature/lock-free-queue
```

### Branch Naming

Use the following format:

```text
feature/<your-concept>
```

The concept should describe what you are building or exploring.

Examples:

```text
feature/tracing
feature/profiling
feature/goroutine-scheduler
feature/lock-free-queue
feature/custom-allocator
feature/runtime-experiment
```

---

## 2. Create Your Own Concept

You are encouraged to create your own concept rather than only extending existing examples.

A concept can be a small experiment or a complete implementation.

For example:

```text
tracing/
├── main.go
└── ...
```

or:

```text
lock_free_queue/
├── queue.go
├── queue_test.go
└── ...
```

The concept does not need to be large.

A focused experiment is enough.

---

## 3. Update README.md

For a new topic, **README.md is the primary documentation**.

You should update the project's `README.md` to include your new topic and a short description.

For example:

```markdown
| `tracing` | Distributed tracing and request tracking experiments. |
```

You do **not** need to create an additional `README.md`, `info.md`, or other documentation file for every topic.

Keep the documentation simple and focused.

### Additional Documentation

Additional documentation is only necessary when the topic is complex enough that a short entry in `README.md` is not sufficient.

Simple topics do not need additional documentation.

---

## 4. Keep Examples Focused

Try to keep each topic focused on one main concept.

For example:

```text
channel
channel_buffered
channel_unbuffered
```

are better separated than putting unrelated channel concepts into one large example.

The goal is to make each topic easy to understand and experiment with independently.

---

## 5. Verify Your Changes

Before submitting a Pull Request, make sure your changes compile and the implementation behaves as intended.

If your contribution includes tests, run them:

```bash
go test ./...
```

For benchmarks:

```bash
go test -bench=. ./...
```

If there are no tests, you should still verify that:

- The code compiles successfully.
- The implementation follows the intended logic.
- The example behaves as expected.
- There are no obvious runtime or logical errors.

For example:

```bash
go build ./...
```

can be used to verify that the project builds successfully.

Tests are encouraged when they are useful, but **a test suite is not required for every contribution**.

---

## 6. Commit Messages

Use clear commit messages.

Recommended prefixes include:

```text
feat:     Add a new feature or topic
fix:      Fix an existing implementation
docs:     Update documentation
test:     Add or modify tests
refactor: Refactor existing code
perf:     Improve performance
```

Examples:

```text
feat: add tracing example
```

```text
feat: add goroutine scheduler experiment
```

```text
fix: correct channel example
```

```text
docs: update topic documentation
```

---

## 7. Update AUTHORS.md

Every contributor is encouraged to add their name and email to `AUTHORS.md`.

Use the following format:

```text
Name | Email
```

For example:

```text
Perry Wang | perry@example.com
```

When contributing, add your own entry:

```text
Your Name | your-email@example.com
```

Do not remove or modify another contributor's entry.

The purpose of `AUTHORS.md` is to give proper credit to everyone who contributes to the project.

---

## 8. Pull Request

After completing your feature:

```bash
git add .
git commit -m "feat: add <your-concept>"
```

Push your branch:

```bash
git push origin feature/<your-concept>
```

Then open a Pull Request.

### Pull Request Title

Use a clear title:

```text
feat: add <your-concept>
```

For example:

```text
feat: add tracing example
```

### Pull Request Description

A simple PR description can follow this structure:

```markdown
## What

Add a tracing example.

## Why

Explore how tracing can be implemented in a Go application.

## Changes

- Add tracing example
- Update README.md
- Update AUTHORS.md

## Verification

go test ./...
```

If the contribution does not have tests:

```markdown
## Verification

- go build ./...
- Manually verified the implementation and expected behavior.
```

Keep Pull Requests focused on the concept being introduced or changed.

---

## 9. Respect Existing Contributions

Please respect existing work and contributors.

When contributing:

- Do not delete another contributor's work without discussion.
- Do not overwrite another contributor's topic.
- Avoid unrelated changes in the same Pull Request.
- Discuss duplicated or overlapping topics when necessary.
- Keep examples simple and focused.
- Explain why the implementation exists, not only how it works.

---

## 10. Contribution Workflow

The complete workflow is:

```text
Idea
  ↓
Create your concept
  ↓
Fork and create feature branch
  ↓
Implement
  ↓
Update README.md
  ↓
Update AUTHORS.md
  ↓
Verify
  ↓
Commit
  ↓
Push
  ↓
Open Pull Request
```

For most contributions, the only documentation you need to update is `README.md`.

Additional documentation is optional and should only be added when the concept is complex enough to require it.

Tests are encouraged when applicable. If tests are not present, make sure the code compiles and the implementation logic is correct.

---

## Final Principle

There is no requirement to build something huge.

Start with an idea.

Turn the idea into a concept.

Turn the concept into code.

Keep the documentation simple.

Verify that it works.

Then share it with others.

> **Create your own concept. Implement it. Document it. Share it.**