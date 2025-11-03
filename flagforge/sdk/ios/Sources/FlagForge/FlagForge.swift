import Foundation

public final class FlagForge {
    public static let shared = FlagForge()

    private var clientKey: String?
    private var environment: String?
    private var baseURL: URL = URL(string: "https://api.flagforge.dev")!
    private var cache: [String: Any] = [:]
    private var observers: [(Dictionary<String, Any>) -> Void] = []
    private var etag: String?

    private init() {}

    public func configure(clientKey: String, environment: String, baseURL: URL? = nil) {
        self.clientKey = clientKey
        self.environment = environment
        if let baseURL {
            self.baseURL = baseURL
        }
        fetchFlags()
    }

    public func bool(_ key: String, default defaultValue: Bool) -> Bool {
        cache[key] as? Bool ?? defaultValue
    }

    public func number(_ key: String, default defaultValue: Double) -> Double {
        if let number = cache[key] as? Double {
            return number
        }
        if let stringValue = cache[key] as? String, let double = Double(stringValue) {
            return double
        }
        return defaultValue
    }

    public func string(_ key: String, default defaultValue: String) -> String {
        cache[key] as? String ?? defaultValue
    }

    public func onUpdate(_ callback: @escaping (Dictionary<String, Any>) -> Void) {
        observers.append(callback)
        callback(cache)
    }

    private func fetchFlags() {
        guard let clientKey, let environment else {
            return
        }

        var components = URLComponents(url: baseURL.appendingPathComponent("/v1/flags"), resolvingAgainstBaseURL: false)
        components?.queryItems = [
            URLQueryItem(name: "env", value: environment),
            URLQueryItem(name: "client_key", value: clientKey)
        ]

        guard let url = components?.url else { return }

        var request = URLRequest(url: url)
        if let etag {
            request.addValue(etag, forHTTPHeaderField: "If-None-Match")
        }

        URLSession.shared.dataTask(with: request) { data, response, error in
            guard error == nil, let http = response as? HTTPURLResponse else {
                return
            }

            guard http.statusCode == 200, let data else {
                return
            }

            if let newEtag = http.value(forHTTPHeaderField: "ETag") {
                self.etag = newEtag
            }

            do {
                let payload = try JSONSerialization.jsonObject(with: data, options: [])
                if let dict = payload as? [String: Any], let flags = dict["flags"] as? [[String: Any]] {
                    var snapshot: [String: Any] = [:]
                    for flag in flags {
                        if let key = flag["key"] as? String {
                            snapshot[key] = flag["value_json"] ?? flag
                        }
                    }
                    self.cache = snapshot
                    self.observers.forEach { $0(snapshot) }
                }
            } catch {
                // TODO: persist error for debugging
            }
        }.resume()
    }
}
