# Fyne Components - Design Principles & Guidelines

**Version**: 1.0  
**Date**: February 20, 2026  
**Status**: Living Document

---

## Table of Contents

- [📖 Purpose](#-purpose)
- [🎯 Module Philosophy](#-module-philosophy)
  - [Core Principles](#core-principles)
  - [Non-Goals](#non-goals)
- [🏗️ Architecture Decisions](#️-architecture-decisions)
  - [1. Directory Structure: Hybrid Approach](#1-directory-structure-hybrid-approach)
  - [2. Logging: Injectable Interface Pattern](#2-logging-injectable-interface-pattern)
  - [3. Event Handling: Callbacks vs Interfaces](#3-event-handling-callbacks-vs-interfaces)
  - [4. Configuration: Pattern Selection](#4-configuration-pattern-selection)
- [📐 API Design Guidelines](#-api-design-guidelines)
  - [General Principles](#general-principles)
  - [Public API Checklist](#public-api-checklist)
  - [Naming Conventions](#naming-conventions)
  - [State Access](#state-access)
- [🧪 Component Standards](#-component-standards)
  - [1. Documentation](#1-documentation)
  - [2. Testing](#2-testing)
  - [3. Examples](#3-examples)
  - [4. Dependencies](#4-dependencies)
- [🚧 Development Constraints](#-development-constraints)
  - [❌ FORBIDDEN Actions](#-forbidden-actions)
  - [✅ REQUIRED Actions](#-required-actions)
- [📦 Versioning & Stability](#-versioning--stability)
  - [Version Scheme](#version-scheme)
  - [Pre-1.0 Policy (v0.x)](#pre-10-policy-v0x)
  - [Post-1.0 Policy (v1.x+)](#post-10-policy-v1x)
  - [Stability Roadmap](#stability-roadmap)
- [🔧 Component Migration Guidelines](#-component-migration-guidelines)
  - [Extracting from RTK (or other projects)](#extracting-from-rtk-or-other-projects)
- [🤖 LLM Agent Instructions](#-llm-agent-instructions)
  - [When Working on fyne-components](#when-working-on-fyne-components)
- [📚 References](#-references)
- [📝 Document Evolution](#-document-evolution)
- [📋 Changelog](#-changelog)

---

## 📖 **Purpose**

This document defines the architectural decisions, design principles, and development constraints for the **fyne-components** module. It serves as guidance for:

- **Maintainers**: Ensuring consistency across components
- **Contributors**: Understanding design philosophy and patterns
- **LLM Agents**: Automated development assistance with proper constraints
- **Users**: Understanding the module's evolution and stability guarantees

All decisions documented here were made based on lessons learned from production use in the [Resource Tech Kit](https://github.com/digital-michael/rtk) project.

---

## 🎯 **Module Philosophy**

### Core Principles

1. **Production-Proven**: Components are extracted from real-world use, not theoretical designs
2. **Conservative API**: Start small, grow based on actual user needs, not speculation
3. **Zero Dependencies**: Only stdlib + Fyne - no external dependencies
4. **Reusability**: Components must be generic, not domain-specific
5. **Stability**: v1.0+ APIs are stable; v0.x allows iteration

### Non-Goals

- ❌ Comprehensive widget library (only extract what's proven useful)
- ❌ Experimental features (use v0.x or separate repo for experiments)
- ❌ Domain-specific widgets (project management, accounting, etc.)
- ❌ Framework/abstraction layer (Fyne is the framework)

---

## 🏗️ **Architecture Decisions**

### 1. **Directory Structure: Hybrid Approach**

**Decision**: Use both `pkg/` (public API) and `internal/` (private implementation).

```
fyne-components/
├── pkg/                    # Public, importable by external projects
│   ├── table/             # Component package
│   │   ├── table.go       # Public API (types, constructors, methods)
│   │   ├── config.go      # Public configuration
│   │   ├── logger.go      # Logger interface & adapters
│   │   └── doc.go         # Package documentation
│   └── password/          # Another component
├── internal/              # Private, not importable externally
│   └── table/            # Implementation details
│       ├── state.go      # State management
│       ├── keyboard.go   # Keyboard handling
│       ├── mouse.go      # Mouse handling
│       └── rendering.go  # Rendering logic
├── examples/             # Runnable examples
│   ├── table-demo/
│   └── password-demo/
└── cmd/                  # Tools and utilities
```

**Rationale**:
- **Small public API**: Only 5-10% of code needs to be public
- **Implementation freedom**: Can refactor internal/ without breaking users
- **Clear boundaries**: Forces good encapsulation
- **Testing flexibility**: Can expose internals for testing if needed (via export_test.go)

**Trade-offs**:
- More complex structure than flat pkg/
- Must be disciplined about what goes public

**When to make something public**:
- ✅ Types users construct or configure (Table, Config, Options)
- ✅ Methods users call (SetData, GetSelectedCell, Refresh)
- ✅ Interfaces users implement (Logger for custom logging)
- ❌ Internal state management
- ❌ Event handling implementation
- ❌ Rendering details

---

### 2. **Logging: Injectable Interface Pattern**

**Decision**: Logging is optional, injectable, and silent by default.

**Implementation**:
```go
// pkg/table/logger.go
type Logger interface {
    Debug(msg string, keyvals ...interface{})
    Info(msg string, keyvals ...interface{})
    Warn(msg string, keyvals ...interface{})
    Error(msg string, keyvals ...interface{})
}

// NoopLogger is the default (silent)
type NoopLogger struct{}
func (NoopLogger) Debug(string, ...interface{}) {}
func (NoopLogger) Info(string, ...interface{}) {}
func (NoopLogger) Warn(string, ...interface{}) {}
func (NoopLogger) Error(string, ...interface{}) {}

// pkg/table/config.go
type Config struct {
    // ... other fields
    Logger Logger  // nil uses NoopLogger
}
```

**Usage by consumers**:
```go
// Silent by default
table := table.NewTable(&table.Config{Columns: cols})

// With custom logger
table := table.NewTable(&table.Config{
    Columns: cols,
    Logger: myLogger,  // Implements table.Logger
})
```

**Rationale**:
- No surprise console spam
- Users control logging destination/format
- Compatible with any logging library (slog, logrus, zap)
- Follows stdlib patterns (net/http, database/sql)

**Constraints**:
- ❌ NEVER use global loggers (log.Println, fmt.Println)
- ❌ NEVER require logging configuration
- ✅ Provide StdLogger adapter for stdlib log package
- ✅ Debug logs for development, Info/Warn/Error for production issues

---

### 3. **Event Handling: Callbacks vs Interfaces**

**Decision**: Start with callbacks, add interfaces only if users request.

**Current Approach (v0.x)**:
```go
type Config struct {
    OnCellEdited   func(row int, col string, newValue string, data interface{})
    OnRowSelected  func(row int, data interface{})
    OnKeyPressed   func(key *fyne.KeyEvent) bool  // Return true to prevent default
    // ... other callbacks
}
```

**Future Possibility (v1.0+, if requested)**:
```go
type KeyHandler interface {
    HandleKey(key *fyne.KeyEvent, table *Table) bool
}

type Config struct {
    KeyHandler KeyHandler  // Optional, overrides default behavior
    // ... existing callbacks still supported
}
```

**Rationale**:
- **Callbacks**: Simple, cover 95% of use cases, no interface implementation needed
- **Interfaces**: Advanced users can completely override behavior
- **Start simple**: Add complexity only when proven necessary

**Guidelines**:
- ✅ Callbacks for event notification (OnCellEdited, OnRowSelected)
- ✅ Callbacks that can prevent default behavior return bool
- ❌ Don't expose handler interfaces until v1.0+ and users demonstrate need
- ✅ When adding interfaces, maintain callback compatibility

---

### 4. **Configuration: Pattern Selection**

**Decision**: Use Config struct for migrated components, Builder pattern for new components.

**Config Struct Pattern** (Table, Password - migrated from RTK):
```go
config := table.NewTableConfig("myTable")
config.Columns = []table.ColumnConfig{...}
config.ShowHeaders = true
config.OnCellEdited = func(...) {...}

table := table.NewTable(config)
```

**Builder Pattern** (New components going forward):
```go
widget := newcomponent.New("id",
    newcomponent.WithOption1(value),
    newcomponent.WithOption2(value),
    newcomponent.OnEvent(func() {...}),
)
```

**Rationale**:
- **Config struct**: Familiar to RTK users, works well for complex configuration
- **Builder pattern**: Modern Go idiom, self-documenting, optional parameters
- **Don't break what works**: RTK components already use Config successfully

**When to use which**:
- **Config struct**: Components with >10 configuration options
- **Builder pattern**: Components with <10 options, new components
- **Both**: Can coexist (Builder can construct Config internally)

**Migration consideration**: Config structs can evolve to support builders later without breaking changes.

---

## 📐 **API Design Guidelines**

### General Principles

1. **Conservative Exposure**: Only make public what users need
2. **95% Use Case**: Optimize for common usage, not edge cases
3. **Explicit State**: No hidden state - all state accessible via methods
4. **Backward Compatibility**: v1.0+ doesn't break; v0.x tries to maintain compatibility
5. **Fail Fast**: Panic on programmer errors (nil config), return errors for runtime issues

### Public API Checklist

Before making a type/method public, ask:

- [ ] Do 50%+ of users need this?
- [ ] Is the name clear and unambiguous?
- [ ] Is the signature stable (unlikely to change)?
- [ ] Is it needed for testing/debugging?
- [ ] Does it expose implementation details?

If last question is "yes", keep it internal.

### Naming Conventions

**Packages**: Lowercase, single word (table, password, dialog)

**Types**: 
- Component: `Table`, `PasswordStrengthMeter`
- Config: `Config`, `TableConfig`, `Options`
- Interfaces: `Logger`, `Renderer`, `Handler`

**Methods**:
- Getters: `GetSelectedCell()`, `GetData()`
- Setters: `SetSelectedCell()`, `SetData()`
- Actions: `Refresh()`, `Clear()`, `Focus()`
- Queries: `IsVisible()`, `HasSelection()`

**Callbacks**:
- Pattern: `On<Event>` (OnCellEdited, OnRowSelected)
- Return bool if preventing default behavior

### State Access

**Rule**: All meaningful state must be queryable and settable.

**Required methods for interactive widgets**:
```go
// Query state
func (w *Widget) GetSelection() (row, col int)
func (w *Widget) GetData() []interface{}
func (w *Widget) IsEditing() bool

// Set state
func (w *Widget) SetSelection(row, col int)
func (w *Widget) SetData(data []interface{})
func (w *Widget) ClearSelection()
```

**Rationale**: Users integrate widgets into complex UIs and need state synchronization.

---

## 🧪 **Component Standards**

Every component must have:

### 1. **Documentation**

- [ ] `README.md` - Usage guide with examples
- [ ] `doc.go` - Package-level documentation
- [ ] Godoc comments on all public types/methods/functions
- [ ] At least one working example in `examples/`

**doc.go template**:
```go
// Package component provides [brief description].
//
// [Feature list]
//
// # Example
//
// [Code example]
//
// # Configuration
//
// [Configuration notes]
package component
```

### 2. **Testing**

- [ ] Unit tests for public API: `component_test.go`
- [ ] Tests for internal logic: `internal/component/*_test.go`
- [ ] Minimum 70% coverage: `make coverage`
- [ ] Table tests for edge cases
- [ ] Benchmark tests for performance-critical code

**Test organization**:
```go
func TestComponentCreation(t *testing.T) { ... }
func TestComponentStateManagement(t *testing.T) { ... }
func TestComponentCallbacks(t *testing.T) { ... }
```

### 3. **Examples**

- [ ] Runnable example in `examples/<component>-demo/`
- [ ] README explaining what the example demonstrates
- [ ] Can be run with: `go run examples/<component>-demo/main.go`

### 4. **Dependencies**

- [ ] Only Fyne + stdlib
- [ ] No external Go dependencies
- [ ] No OS-specific code (or clearly documented/optional)

---

## 🚧 **Development Constraints**

### For Maintainers and LLM Agents

#### ❌ **FORBIDDEN Actions**

1. **Adding Dependencies**
   - Never add go.mod dependencies without explicit discussion
   - Exception: Fyne updates for bug fixes
   - Rationale: Keep module lightweight and build fast

2. **Breaking Public API (v1.0+)**
   - Cannot remove public types/methods/functions
   - Cannot change signatures of public API
   - Cannot change behavior of documented features
   - Must use deprecation cycle (mark deprecated, wait one version, then remove)

3. **Exposing Internals**
   - Cannot make internal/ packages public
   - Cannot add public methods that leak implementation details
   - Cannot return internal types from public methods

4. **Global State**
   - No package-level vars for state
   - No init() functions that mutate global state
   - No global loggers or configuration

5. **Domain-Specific Code**
   - Components must be generic/reusable
   - No business logic (validation, workflows, etc.)
   - No RTK-specific features

#### ✅ **REQUIRED Actions**

1. **Before Adding Features**
   - Demonstrate need with real use case
   - Consider if 50%+ of users need it
   - Start with callback, defer interface until requested
   - Add tests for new functionality

2. **Before Making Changes**
   - Check if it breaks existing tests
   - Update documentation if behavior changes
   - Consider backward compatibility
   - Add migration guide for breaking changes (v0.x → v1.0)

3. **Code Quality**
   - Run `make fmt` before committing
   - Run `make lint` and fix all issues
   - Run `make test` and ensure all pass
   - Maintain >70% test coverage

4. **Documentation**
   - Update godoc for changed methods
   - Update README if API changes
   - Update examples if usage pattern changes
   - Add entry to CHANGELOG.md

---

## 📦 **Versioning & Stability**

### Version Scheme

**Semantic Versioning**: `MAJOR.MINOR.PATCH`

- **MAJOR**: Breaking changes (v1.0.0 → v2.0.0)
- **MINOR**: New features, backward compatible (v1.0.0 → v1.1.0)
- **PATCH**: Bug fixes only (v1.0.0 → v1.0.1)

### Pre-1.0 Policy (v0.x)

**Status**: API is evolving, breaking changes allowed

**Guidelines**:
- MINOR bumps (v0.1.0 → v0.2.0) may include breaking changes
- Document breaking changes in CHANGELOG.md
- Try to maintain compatibility where reasonable
- Use deprecation warnings where possible

**Users should**:
- Pin to specific v0.x.y versions (not v0.x)
- Expect migration when updating minor versions
- Report API issues early for v1.0 consideration

### Post-1.0 Policy (v1.x+)

**Status**: API is stable, no breaking changes

**Guarantees**:
- Existing code continues to work
- New features added without breaking existing
- Deprecations given one full minor version notice
- Bug fixes don't change documented behavior

**Breaking changes**:
- Require MAJOR version bump (v1.x → v2.0)
- Provide migration guide
- Consider v2 as separate module for Go modules compatibility

### Stability Roadmap

- **v0.1.0-alpha**: Initial structure, no stability guarantees
- **v0.1.0-beta**: Feature complete, API mostly stable, testing phase
- **v0.1.0**: First stable v0.x release, basic stability guarantees
- **v1.0.0**: Production ready, full API stability guarantees

**Criteria for v1.0.0**:
- [ ] Used in production for 6+ months
- [ ] 5+ different usage sites
- [ ] API hasn't changed for 3+ months
- [ ] >80% test coverage
- [ ] Complete documentation
- [ ] No known critical bugs

---

## 🔧 **Component Migration Guidelines**

### Extracting from RTK (or other projects)

When migrating components to fyne-components:

1. **Dependency Check**
   - [ ] Verify zero non-Fyne dependencies
   - [ ] List all internal imports
   - [ ] Plan how to eliminate internal dependencies

2. **Public API Design**
   - [ ] Identify what users actually use (not what exists)
   - [ ] Move only necessary API to pkg/
   - [ ] Keep implementation in internal/
   - [ ] Design for 95% use case

3. **Testing**
   - [ ] Port existing tests
   - [ ] Add tests for missing coverage
   - [ ] Verify tests pass standalone (no RTK context)

4. **Documentation**
   - [ ] Write component README
   - [ ] Add godoc to all public API
   - [ ] Create example
   - [ ] Do NOT import private lessons learned from source project

5. **Integration Testing**
   - [ ] Test component in source project with external import
   - [ ] Verify no functionality loss
   - [ ] Fix any integration issues

---

## 🤖 **LLM Agent Instructions**

### When Working on fyne-components

**HIGH PRIORITY RULES**:

1. **Read this file first** before making changes
2. **Check existing patterns** in other components before inventing new ones
3. **Ask before adding dependencies** - even Fyne sub-packages
4. **Test after every change** - run `make test` frequently
5. **Update docs** when changing public API

**When asked to add a feature**:

1. Check if similar feature exists in another component
2. Propose minimal API first (callbacks, not interfaces)
3. Ask: "Do 50%+ of users need this?"
4. Implement in internal/, expose minimally in pkg/
5. Add tests and documentation

**When asked to fix a bug**:

1. Write a failing test first
2. Fix in appropriate layer (internal/ vs pkg/)
3. Verify existing tests still pass
4. Check if fix affects public API (document if so)

**When refactoring**:

1. NEVER change public API without explicit approval
2. Keep internal/ changes internal
3. Maintain backward compatibility
4. Update tests if behavior changes

**Red Flags (Stop and ask)**:

- 🚨 Changing public method signatures
- 🚨 Adding go.mod dependencies
- 🚨 Creating global state
- 🚨 Exposing internal types
- 🚨 Breaking existing tests
- 🚨 Domain-specific features

---

## 📚 **References**

- **Source Project**: [Resource Tech Kit](https://github.com/digital-michael/rtk)
- **UI Framework**: [Fyne v2](https://fyne.io/)
- **Similar Projects**: 
  - [fyne-x](https://github.com/fyne-io/fyne-x) - Official experimental widgets
  - [fynewidgets](https://github.com/andydotxyz/fynewidgets) - Community widgets

---

## 📝 **Document Evolution**

This document is a living guide and should be updated when:

- Major architectural decisions are made
- New patterns are established
- Constraints are added or removed
- Version policy changes

**Update Process**:
1. Propose change in issue or PR
2. Discuss with maintainers
3. Update document with rationale
4. Bump version number at top
5. Add changelog entry below

---

## 📋 **Changelog**

### Version 1.0 (February 20, 2026)
- Initial version based on RTK extraction decisions
- Established hybrid pkg/internal structure
- Defined logging interface pattern
- Established callback-first approach
- Defined Config vs Builder pattern usage
- Set v0.x vs v1.0+ stability guarantees

---

**Questions or Suggestions?** Open an issue or discussion on GitHub.
