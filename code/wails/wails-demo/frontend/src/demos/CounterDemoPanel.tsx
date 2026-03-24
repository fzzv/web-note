import {useEffect, useState} from 'react';
import {Decrement, GetCount, Increment} from '../../wailsjs/go/main/CounterDemo';
import type {Messages} from '../i18n';

type CounterDemoPanelProps = {
    messages: Messages;
};

function CounterDemoPanel({messages}: CounterDemoPanelProps) {
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
        return <div className="status-card">{messages.counterPanel.loading}</div>;
    }

    return (
        <section className="demo-panel">
            <div className="panel-head">
                <p className="panel-kicker">{messages.counterPanel.kicker}</p>
                <h3>{messages.counterPanel.title}</h3>
                <p className="panel-copy">{messages.counterPanel.description}</p>
            </div>

            <div className="counter-display">
                <span className="counter-label">{messages.counterPanel.countLabel}</span>
                <strong>{count}</strong>
            </div>

            <div className="action-row">
                <button className="action-button secondary" onClick={() => void changeCount(Decrement)} type="button">
                    {messages.counterPanel.decrement}
                </button>
                <button className="action-button primary" onClick={() => void changeCount(Increment)} type="button">
                    {messages.counterPanel.increment}
                </button>
            </div>

            <dl className="meta-grid">
                <div>
                    <dt>{messages.counterPanel.startupLabel}</dt>
                    <dd>{messages.counterPanel.startupDescription}</dd>
                </div>
                <div>
                    <dt>{messages.counterPanel.beforeCloseLabel}</dt>
                    <dd>{messages.counterPanel.beforeCloseDescription}</dd>
                </div>
                <div>
                    <dt>{messages.counterPanel.shutdownLabel}</dt>
                    <dd>{messages.counterPanel.shutdownDescription}</dd>
                </div>
            </dl>

            {error ? <p className="feedback error">{error}</p> : null}
        </section>
    );
}

export {CounterDemoPanel};
