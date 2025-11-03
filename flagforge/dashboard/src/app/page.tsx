import Link from 'next/link';
import { FlagTable } from '../components/FlagTable';

export default function HomePage() {
  const projectId = 'demo-project';
  const environment = 'dev';

  return (
    <main className="container">
      <header style={{ marginBottom: '32px' }}>
        <h1 style={{ fontSize: '32px', fontWeight: 600, marginBottom: '8px' }}>FlagForge Dashboard</h1>
        <p style={{ color: '#475569', margin: 0 }}>
          Manage feature flags and remote configuration. This scaffold renders demo data and fetches from the API when available.
        </p>
      </header>

      <section className="card">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
          <div>
            <h2 style={{ margin: 0, fontSize: '20px' }}>Flags</h2>
            <p style={{ margin: '4px 0 0', color: '#64748b' }}>
              Project <code>{projectId}</code> / Environment <code>{environment}</code>
            </p>
          </div>
          <Link href="/" style={{ fontSize: '14px', fontWeight: 600 }}>
            New Flag
          </Link>
        </div>
        <FlagTable projectId={projectId} environment={environment} />
      </section>
    </main>
  );
}
