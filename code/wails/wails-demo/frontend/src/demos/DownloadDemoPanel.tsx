import {useEffect, useState} from 'react';
import {GetState, Reset, StartDownload} from '../../wailsjs/go/main/DownloadDemo';
import {EventsOn} from '../../wailsjs/runtime/runtime';
import type {Messages} from '../i18n';

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

type DownloadDemoPanelProps = {
    messages: Messages;
};

const idleState: DownloadState = {
    status: 'idle',
    message: '',
    url: '',
    fileName: '',
    destination: '',
    downloadedBytes: 0,
    totalBytes: 0,
    progress: 0,
};

function DownloadDemoPanel({messages}: DownloadDemoPanelProps) {
    const [downloadURL, setDownloadURL] = useState(defaultDownloadURL);
    const [state, setState] = useState<DownloadState>({
        ...idleState,
        message: messages.downloadPanel.loading,
    });
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
        return <div className="status-card">{messages.downloadPanel.loading}</div>;
    }

    const inProgress = state.status === 'starting' || state.status === 'downloading';
    const progress = state.status === 'completed' ? 100 : clampProgress(state.progress);
    const statusLabel = messages.downloadPanel.statusLabels[state.status as keyof typeof messages.downloadPanel.statusLabels] ?? state.status;

    return (
        <section className="demo-panel">
            <div className="panel-head">
                <p className="panel-kicker">{messages.downloadPanel.kicker}</p>
                <h3>{messages.downloadPanel.title}</h3>
                <p className="panel-copy" dangerouslySetInnerHTML={{__html: messages.downloadPanel.description}}/>
            </div>

            <div className="input-card">
                <label className="field-label">
                    {messages.downloadPanel.urlLabel}
                    <input
                        onChange={(event) => setDownloadURL(event.target.value)}
                        placeholder={messages.downloadPanel.urlPlaceholder}
                        type="url"
                        value={downloadURL}
                    />
                </label>
                <p className="field-note" dangerouslySetInnerHTML={{__html: messages.downloadPanel.note}}/>
                <div className="action-row">
                    <button className="action-button primary" disabled={inProgress} onClick={() => void handleStart()} type="button">
                        {inProgress ? messages.downloadPanel.downloading : messages.downloadPanel.start}
                    </button>
                    <button className="action-button ghost" disabled={inProgress} onClick={() => void handleReset()} type="button">
                        {messages.downloadPanel.reset}
                    </button>
                </div>
            </div>

            <div className="progress-card">
                <div className="status-row">
                    <span className={`status-pill ${state.status}`}>{statusLabel}</span>
                    <strong className="progress-value">{progress.toFixed(0)}%</strong>
                </div>

                <div className="progress-track">
                    <div className="progress-fill" style={{width: `${progress}%`}}/>
                </div>

                <p className="progress-message">{state.message}</p>
            </div>

            <dl className="meta-grid">
                <div>
                    <dt>{messages.downloadPanel.fileLabel}</dt>
                    <dd>{state.fileName || messages.downloadPanel.waitingFile}</dd>
                </div>
                <div>
                    <dt>{messages.downloadPanel.transferredLabel}</dt>
                    <dd>{formatTransferred(state.downloadedBytes, state.totalBytes)}</dd>
                </div>
                <div>
                    <dt>{messages.downloadPanel.savedToLabel}</dt>
                    <dd>{state.destination || messages.downloadPanel.savedToPlaceholder}</dd>
                </div>
            </dl>

            <ul className="hint-list">
                {messages.downloadPanel.hints.map((hint) => (
                    <li key={hint}>{hint}</li>
                ))}
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
