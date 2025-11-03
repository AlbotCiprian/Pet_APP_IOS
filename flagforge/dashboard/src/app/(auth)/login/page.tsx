export default function LoginPage() {
  return (
    <main style={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', backgroundColor: '#f8fafc', padding: '24px' }}>
      <div style={{ width: '100%', maxWidth: '360px', borderRadius: '12px', border: '1px solid #e2e8f0', background: '#ffffff', padding: '24px', boxShadow: '0 12px 32px rgba(15,23,42,0.08)' }}>
        <header style={{ textAlign: 'center', marginBottom: '24px' }}>
          <h1 style={{ fontSize: '24px', fontWeight: 600, marginBottom: '8px', color: '#0f172a' }}>Sign in to FlagForge</h1>
          <p style={{ fontSize: '14px', color: '#64748b', margin: 0 }}>Demo credentials are not required yet.</p>
        </header>
        <form style={{ display: 'grid', gap: '16px' }}>
          <div style={{ display: 'grid', gap: '4px' }}>
            <label style={{ fontSize: '14px', fontWeight: 600, color: '#0f172a' }} htmlFor="email">
              Email
            </label>
            <input id="email" type="email" placeholder="you@example.com" style={{ width: '100%', borderRadius: '8px', border: '1px solid #cbd5f5', padding: '10px 12px', fontSize: '14px' }} />
          </div>
          <div style={{ display: 'grid', gap: '4px' }}>
            <label style={{ fontSize: '14px', fontWeight: 600, color: '#0f172a' }} htmlFor="password">
              Password
            </label>
            <input id="password" type="password" placeholder="••••••••" style={{ width: '100%', borderRadius: '8px', border: '1px solid #cbd5f5', padding: '10px 12px', fontSize: '14px' }} />
          </div>
          <button type="button" style={{ width: '100%', borderRadius: '8px', background: '#4f46e5', color: '#ffffff', fontWeight: 600, fontSize: '14px', padding: '10px 12px', border: 'none', cursor: 'pointer' }}>
            Continue
          </button>
        </form>
      </div>
    </main>
  );
}
