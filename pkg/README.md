# pkg/

This directory contains the public library code that can be imported by external projects.

Place your reusable, publicly exported packages here. Each subdirectory should represent a logical package that external users can import.

Example structure:
```
pkg/
├── widgets/       # Custom Fyne widgets
├── layouts/       # Custom layouts
└── themes/        # Custom themes
```

Import example:
```go
import "github.com/digital-michael/fyne-components/pkg/widgets"
```
