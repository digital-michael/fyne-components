# internal/

This directory contains private application code that cannot be imported by external projects.

Go enforces that code in `internal/` directories cannot be imported by projects outside the parent tree. This is ideal for:
- Private utilities
- Implementation details
- Helper functions not meant for public consumption

Example structure:
```
internal/
├── utils/         # Private utility functions
└── helpers/       # Internal helper code
```
