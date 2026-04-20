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

### 1. Project Management
- `kodia new <project-name>`: Scaffolds a new professional fullstack application (Go + SvelteKit).
- `kodia serve [--watch]`: Starts the developmental server with automatic hot-reload for the backend.
- `kodia tinker`: Opens an interactive REPL shell to test your business logic and database queries in real-time.

### 2. Artisanal Scaffolding (`generate`)
Kodia features a powerful generation engine to keep your arsitektur clean and consistent.

```bash
# Generate a full CRUD resource
kodia generate crud users --fields=name:string,email:string,role:enum

# Individual components
kodia generate model Post title:string body:text author:references:User
kodia generate migration AddStatusToUsers
kodia generate event UserCreated
kodia generate listener SendWelcomeEmail --event=UserCreated
kodia generate policy UserPolicy
kodia generate test Feature/AuthenticationTest --type=feature
```

### 3. Plugin Ecosystem (`plugin`)
Manage official and community-built extensions to keep your core framework lightweight.

```bash
# Install an external plugin
kodia plugin install payment
```

### 4. Database Management (`db`)
- `kodia db:migrate`: Run pending migrations.
- `kodia db:rollback`: Roll back the last migration.
- `kodia db:seed`: Seed the database with fake data using artisanal seeders.

---

## 🏗️ Architecture Focus

Unlike generic generators, the Kodia CLI produces **Modular Architecture** compliant code:
- **Separation of Concerns**: Logic is divided into Handlers, Services, and Repositories.
- **Type Safety**: Automatically generates strongly-typed DTOs and GORM models.
- **Professional Patterns**: Implements Dependency Injection in every generated component.

---

## 📜 License

Kodia is open-source software licensed under the [MIT license](https://opensource.org/licenses/MIT).

---

<p align="center">
  Built with 🐨 by the Kodia Colony.
</p>
