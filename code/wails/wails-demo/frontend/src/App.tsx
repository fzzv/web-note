import {useEffect, useState} from 'react';
import './App.css';
import {Decrement, GetCount, Increment} from "../wailsjs/go/main/App";

function App() {
    const [count, setCount] = useState(0);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        let cancelled = false;

        GetCount()
            .then((value) => {
                if (!cancelled) {
                    setCount(value);
                }
            })
            .catch((err) => {
                if (!cancelled) {
                    setError(err instanceof Error ? err.message : String(err));
                }
            })
            .finally(() => {
                if (!cancelled) {
                    setLoading(false);
                }
            });

        return () => {
            cancelled = true;
        };
    }, []);

    const changeCount = async (action: () => Promise<number>) => {
        setError('');

        try {
            const nextCount = await action();
            setCount(nextCount);
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        }
    };

    if (loading) {
        return (
            <div className="page-shell">
                <div className="status-card">Loading counter state...</div>
            </div>
        );
    }

    return (
        <div className="page-shell">
            <section className="counter-card">
                <div className="counter-header">
                    <p className="eyebrow">Lifecycle Demo</p>
                    <h1>Persistent Counter</h1>
                    <p className="subtitle">
                        The count is loaded during startup, confirmed before close, and written back on shutdown.
                    </p>
                </div>

                <div className="counter-display">
                    <span className="counter-label">Current count</span>
                    <strong>{count}</strong>
                </div>

                <div className="actions">
                    <button className="action-button secondary" onClick={() => void changeCount(Decrement)}>
                        Decrement
                    </button>
                    <button className="action-button primary" onClick={() => void changeCount(Increment)}>
                        Increment
                    </button>
                </div>

                <div className="notes">
                    <p>Startup reads the previous count from a local JSON file.</p>
                    <p>Before closing, Wails shows a confirmation dialog.</p>
                    <p>Shutdown persists the latest count automatically.</p>
                </div>

                {error ? <p className="feedback error">{error}</p> : null}
            </section>
        </div>
    );
}

export default App;
