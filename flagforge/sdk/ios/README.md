# FlagForge iOS SDK (Preview)

This Swift Package provides a lightweight client for consuming FlagForge feature flags.

```swift
import FlagForge

FlagForge.shared.configure(clientKey: "client-key", environment: "prod")
let isEnabled = FlagForge.shared.bool("new-home", default: false)

// Optionally provide a custom API base during development
FlagForge.shared.configure(
    clientKey: "client-key",
    environment: "prod",
    baseURL: URL(string: "http://localhost:8080")!
)
```

> **Note:** Network persistence, offline caching, and secure storage are left as TODOs for implementers.
