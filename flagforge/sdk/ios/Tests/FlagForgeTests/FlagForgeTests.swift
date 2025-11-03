import XCTest
@testable import FlagForge

final class FlagForgeTests: XCTestCase {
    func testDefaultsReturnedWhenCacheEmpty() {
        let sdk = FlagForge.shared
        sdk.configure(clientKey: "demo", environment: "dev")
        XCTAssertEqual(sdk.bool("missing", default: true), true)
        XCTAssertEqual(sdk.number("missing", default: 42), 42)
        XCTAssertEqual(sdk.string("missing", default: "fallback"), "fallback")
    }
}
