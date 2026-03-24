import {useEffect, useState} from 'react';
import {Decrement, GetCount, Increment} from '../../wailsjs/go/main/CounterDemo';

function CounterDemoPanel() {
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
        return <div className="status-card">Loading counter state...</div>;
    }

    return (
        <section className="demo-panel">
            <div className="panel-head">
                <p className="panel-kicker">Isolated module</p>
                <h3>Counter demo stays intact.</h3>
                <p className="panel-copy">
                    The counter logic now lives in its own backend file, but it still restores state during startup and writes JSON on shutdown.
                </p>
            </div>

            <div className="counter-display">
                <span className="counter-label">Current count</span>
                <strong>{count}</strong>
            </div>

            <div className="action-row">
                <button className="action-button secondary" onClick={() => void changeCount(Decrement)} type="button">
                    Decrement
                </button>
                <button className="action-button primary" onClick={() => void changeCount(Increment)} type="button">
                    Increment
                </button>
            </div>

            <dl className="meta-grid">
                <div>
                    <dt>Startup</dt>
                    <dd>Reads the saved count from a local JSON file before the React UI renders.</dd>
                </div>
                <div>
                    <dt>Before Close</dt>
                    <dd>Shows a confirmation dialog so lifecycle hooks are visible while testing the demo.</dd>
                </div>
                <div>
                    <dt>Shutdown</dt>
                    <dd>Persists the current count automatically, so reopening the app resumes the previous state.</dd>
                </div>
            </dl>

            {error ? <p className="feedback error">{error}</p> : null}
        </section>
    );
}

export {CounterDemoPanel};
