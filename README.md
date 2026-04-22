# Kodia CLI 🐨⚡

The high-velocity command-line companion for the **Kodia Framework**. 

Kodia CLI is designed to eliminate boilerplate and accelerate your fullstack development workflow. It handles project scaffolding, artisanal code generation, dependency management, and real-time development utilities.

---

## 🛠️ Installation

The easiest way to install the Kodia CLI is via `go install`:

```bash
go install github.com/kodia-studio/cli/kodia@latest
```

Ensuring your `$GOPATH/bin` is in your system's `PATH`.

Alternatively, you can build from source:

```bash
# Clone the repository
git clone https://github.com/kodia-studio/kodia.git

# Build and install locally
cd cli
go build -o kodia ./kodia/main.go
mv kodia /usr/local/bin/ # Or add to your PATH
```

---

## 🚀 Core Commands

### 1. Project Management & Testing
- `kodia new <project-name>`: Scaffolds a new professional fullstack application (Go + SvelteKit).
- `kodia test`: Run all backend and frontend tests in the project.
- `kodia test:coverage [--html]`: Generate and view a comprehensive coverage report.
- `kodia shell`: Opens an interactive Go REPL shell to test your business logic in real-time.
- `kodia route:list`: Display a beautiful map of all registered API routes.

### 2. Pro Scaffolding & Security
```bash
# Generate a full vertical slice module
kodia make:module Product

# Secure your application by generating secure cryptographic keys
kodia key:generate
```

### 3. Build & Deployment (`build`)
Kodia features a world-class build pipeline that bundles your entire stack into a single file.

```bash
# Compile backend and frontend into a single production-ready binary
kodia build
```

### 4. Database & Environment (`migrate`, `env`)
- `kodia migrate`: Run all pending database migrations.
- `kodia env:set KEY=VALUE`: Securely update your application configuration.
- `kodia env:list`: View all environment variables (with sensitive data masked).
- `kodia db:seed`: Populate the database with dummy data.

---

## 🏗️ Architecture Focus

Kodia CLI produces **Elite Modular Architecture** compliant code:
- **Separation of Concerns**: Logic is divided into Handlers, Services, and Repositories.
- **Type Safety**: Strongly-typed DTOs and GORM models.
- **Pro Auto-Wiring**: Automatically registers your new modules in `main.go` and `router.go`.

---

## 📜 License

Kodia is open-source software licensed under the [MIT license](https://opensource.org/licenses/MIT).

---

<p align="center">
  Built with 🐨 by the Kodia Colony.
</p>
