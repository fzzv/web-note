import {useEffect, useState} from 'react';
import {GetState, Reset, StartDownload} from '../../wailsjs/go/main/DownloadDemo';
import {EventsOn} from '../../wailsjs/runtime/runtime';

const downloadProgressEventName = 'demo:download:progress';
const defaultDownloadURL = 'https://raw.githubusercontent.com/wailsapp/wails/master/README.md';

type DownloadState = {
    status: string;
    message: string;
    url: string;
    fileName: string;
    destination: string;
    downloadedBytes: number;
    totalBytes: number;
    progress: number;
};

const idleState: DownloadState = {
    status: 'idle',
    message: 'Enter a URL to stream a file and watch backend progress events in real time.',
    url: '',
    fileName: '',
    destination: '',
    downloadedBytes: 0,
    totalBytes: 0,
    progress: 0,
};

function DownloadDemoPanel() {
    const [downloadURL, setDownloadURL] = useState(defaultDownloadURL);
    const [state, setState] = useState<DownloadState>(idleState);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    useEffect(() => {
        let cancelled = false;

        const unsubscribe = EventsOn(downloadProgressEventName, (...payload) => {
            if (cancelled || payload.length === 0) {
                return;
            }

            setState(normalizeState(payload[0]));
        });

        GetState()
            .then((value) => {
                if (!cancelled) {
                    const nextState = normalizeState(value);
                    setState(nextState);
                    if (nextState.url) {
                        setDownloadURL(nextState.url);
                    }
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
            unsubscribe();
        };
    }, []);

    const handleStart = async () => {
        setError('');

        try {
            const nextState = await StartDownload(downloadURL);
            setState(normalizeState(nextState));
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        }
    };

    const handleReset = async () => {
        setError('');

        try {
            const nextState = await Reset();
            setState(normalizeState(nextState));
        } catch (err) {
            setError(err instanceof Error ? err.message : String(err));
        }
    };

    if (loading) {
        return <div className="status-card">Loading download demo...</div>;
    }

    const inProgress = state.status === 'starting' || state.status === 'downloading';
    const progress = state.status === 'completed' ? 100 : clampProgress(state.progress);

    return (
        <section className="demo-panel">
            <div className="panel-head">
                <p className="panel-kicker">runtime.EventsEmit</p>
                <h3>Watch backend download progress arrive as events.</h3>
                <p className="panel-copy">
                    Starting a download returns immediately. The Go backend keeps streaming the response, writes it to disk, and emits progress objects that React subscribes to with <code>EventsOn</code>.
                </p>
            </div>

            <div className="input-card">
                <label className="field-label">
                    Download URL
                    <input
                        onChange={(event) => setDownloadURL(event.target.value)}
                        placeholder="https://example.com/file.zip"
                        type="url"
                        value={downloadURL}
                    />
                </label>
                <p className="field-note">
                    The file is saved into your local <code>Downloads/wails-demo</code> folder. Use any reachable HTTP or HTTPS file URL.
                </p>
                <div className="action-row">
                    <button className="action-button primary" disabled={inProgress} onClick={() => void handleStart()} type="button">
                        {inProgress ? 'Downloading...' : 'Start download'}
                    </button>
                    <button className="action-button ghost" disabled={inProgress} onClick={() => void handleReset()} type="button">
                        Reset status
                    </button>
                </div>
            </div>

            <div className="progress-card">
                <div className="status-row">
                    <span className={`status-pill ${state.status}`}>{state.status || 'idle'}</span>
                    <strong className="progress-value">{progress.toFixed(0)}%</strong>
                </div>

                <div className="progress-track">
                    <div className="progress-fill" style={{width: `${progress}%`}}/>
                </div>

                <p className="progress-message">{state.message}</p>
            </div>

            <dl className="meta-grid">
                <div>
                    <dt>File</dt>
                    <dd>{state.fileName || 'Waiting for a new download'}</dd>
                </div>
                <div>
                    <dt>Transferred</dt>
                    <dd>{formatTransferred(state.downloadedBytes, state.totalBytes)}</dd>
                </div>
                <div>
                    <dt>Saved To</dt>
                    <dd>{state.destination || 'Downloads/wails-demo'}</dd>
                </div>
            </dl>

            <ul className="hint-list">
                <li>Go starts the request in a goroutine, so the UI thread is never blocked.</li>
                <li>Each event carries a full <code>DownloadState</code> payload that React can render directly.</li>
                <li>When the request finishes or fails, the backend emits one final terminal state.</li>
            </ul>

            {error ? <p className="feedback error">{error}</p> : null}
        </section>
    );
}

function normalizeState(value: Partial<DownloadState> | undefined): DownloadState {
    return {
        status: value?.status ?? idleState.status,
        message: value?.message ?? idleState.message,
        url: value?.url ?? '',
        fileName: value?.fileName ?? '',
        destination: value?.destination ?? '',
        downloadedBytes: value?.downloadedBytes ?? 0,
        totalBytes: value?.totalBytes ?? 0,
        progress: value?.progress ?? 0,
    };
}

function clampProgress(progress: number): number {
    if (Number.isNaN(progress)) {
        return 0;
    }

    return Math.max(0, Math.min(progress, 100));
}

function formatTransferred(downloadedBytes: number, totalBytes: number): string {
    if (totalBytes <= 0) {
        return formatBytes(downloadedBytes);
    }

    return `${formatBytes(downloadedBytes)} / ${formatBytes(totalBytes)}`;
}

function formatBytes(size: number): string {
    if (size < 1024) {
        return `${size} B`;
    }

    const units = ['KB', 'MB', 'GB', 'TB'];
    let value = size / 1024;
    let unitIndex = 0;

    while (value >= 1024 && unitIndex < units.length - 1) {
        value /= 1024;
        unitIndex += 1;
    }

    return `${value.toFixed(1)} ${units[unitIndex]}`;
}

export {DownloadDemoPanel};
