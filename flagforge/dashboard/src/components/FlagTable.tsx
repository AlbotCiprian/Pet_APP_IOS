'use client';

import { useEffect, useState } from 'react';

export type Flag = {
  flag_id: string;
  environment_id: string;
  value_json: unknown;
  rollout: number;
  rules_json: unknown;
  updated_at: string;
  updated_by: string;
  key?: string;
  type?: string;
};

interface Props {
  projectId: string;
  environment: string;
}

export function FlagTable({ projectId, environment }: Props) {
  const [flags, setFlags] = useState<Flag[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const controller = new AbortController();
    const apiBase = process.env.NEXT_PUBLIC_API_BASE_URL ?? 'http://localhost:8080';
    setLoading(true);
    fetch(`${apiBase}/v1/flags?project_id=${projectId}&env=${environment}`, {
      signal: controller.signal
    })
      .then(async (res) => {
        if (!res.ok) {
          throw new Error('Failed to load flags');
        }
        const data = await res.json();
        setFlags(Array.isArray(data.flags) ? data.flags : []);
      })
      .catch((err) => {
        if (err.name !== 'AbortError') {
          setError(err.message);
        }
      })
      .finally(() => setLoading(false));

    return () => controller.abort();
  }, [projectId, environment]);

  if (loading) {
    return <p style={{ color: '#64748b', fontSize: '14px' }}>Loading flags…</p>;
  }

  if (error) {
    return (
      <div style={{ border: '1px solid #fecaca', backgroundColor: '#fee2e2', padding: '12px', color: '#b91c1c', borderRadius: '8px', fontSize: '14px' }}>
        {error}
      </div>
    );
  }

  if (!flags.length) {
    return <p style={{ color: '#64748b', fontSize: '14px' }}>No flags found yet.</p>;
  }

  return (
    <table>
      <thead>
        <tr>
          <th>Flag Key</th>
          <th>Type</th>
          <th>Environment</th>
          <th>Value</th>
          <th>Rollout</th>
          <th>Last Updated</th>
        </tr>
      </thead>
      <tbody>
        {flags.map((flag) => (
          <tr key={`${flag.flag_id}-${flag.environment_id}`}>
            <td style={{ fontFamily: 'monospace', fontSize: '12px' }}>{flag.key ?? flag.flag_id}</td>
            <td>{flag.type ?? 'unknown'}</td>
            <td>{environment}</td>
            <td style={{ fontFamily: 'monospace', fontSize: '12px' }}>
              {typeof flag.value_json === 'string'
                ? flag.value_json
                : JSON.stringify(flag.value_json)}
            </td>
            <td>{flag.rollout}%</td>
            <td>{flag.updated_at ? new Date(flag.updated_at).toLocaleString() : '—'}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
