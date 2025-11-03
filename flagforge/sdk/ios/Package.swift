// swift-tools-version:5.10
import PackageDescription

let package = Package(
    name: "FlagForge",
    platforms: [
        .iOS(.v13)
    ],
    products: [
        .library(
            name: "FlagForge",
            targets: ["FlagForge"]
        )
    ],
    dependencies: [],
    targets: [
        .target(
            name: "FlagForge",
            dependencies: []
        ),
        .testTarget(
            name: "FlagForgeTests",
            dependencies: ["FlagForge"]
        )
    ]
)
