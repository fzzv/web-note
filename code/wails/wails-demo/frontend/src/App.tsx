import {useEffect, useMemo, useState} from 'react';
import './App.css';
import {CounterDemoPanel} from './demos/CounterDemoPanel';
import {DownloadDemoPanel} from './demos/DownloadDemoPanel';
import {GetLanguage, SetLanguage} from '../wailsjs/go/main/App';
import {getMessages, isLanguage, languageOptions, normalizeLanguage, type Language} from './i18n';

type DemoId = 'counter' | 'download';

type DemoDefinition = {
    id: DemoId;
    badge: string;
    title: string;
    summary: string;
    focus: string;
};

const languageStorageKey = 'wails-demo.language';

function App() {
    const [activeDemo, setActiveDemo] = useState<DemoId>('counter');
    const [language, setLanguage] = useState<Language>(() => getStoredLanguage() ?? 'zh-CN');

    const messages = useMemo(() => getMessages(language), [language]);
    const demos: DemoDefinition[] = useMemo(() => ([
        {
            id: 'counter',
            badge: messages.demos.counter.badge,
            title: messages.demos.counter.title,
            summary: messages.demos.counter.summary,
            focus: messages.demos.counter.focus,
        },
        {
            id: 'download',
            badge: messages.demos.download.badge,
            title: messages.demos.download.title,
            summary: messages.demos.download.summary,
            focus: messages.demos.download.focus,
        },
    ]), [messages]);

    const currentDemo = demos.find((demo) => demo.id === activeDemo) ?? demos[0];

    useEffect(() => {
        let cancelled = false;
        const storedLanguage = getStoredLanguage();

        GetLanguage()
            .then(async (backendLanguage) => {
                const nextLanguage = storedLanguage ?? normalizeLanguage(backendLanguage);
                const syncedLanguage = normalizeLanguage(await SetLanguage(nextLanguage));
                if (!cancelled) {
                    applyLanguage(syncedLanguage, setLanguage);
                }
            })
            .catch(async () => {
                const fallbackLanguage = storedLanguage ?? language;

                try {
                    const syncedLanguage = normalizeLanguage(await SetLanguage(fallbackLanguage));
                    if (!cancelled) {
                        applyLanguage(syncedLanguage, setLanguage);
                    }
                } catch {
                    if (!cancelled) {
                        applyLanguage(fallbackLanguage, setLanguage);
                    }
                }
            });

        return () => {
            cancelled = true;
        };
    }, []);

    const handleLanguageChange = async (nextLanguage: Language) => {
        applyLanguage(nextLanguage, setLanguage);

        try {
            const syncedLanguage = normalizeLanguage(await SetLanguage(nextLanguage));
            applyLanguage(syncedLanguage, setLanguage);
        } catch {
            applyLanguage(nextLanguage, setLanguage);
        }
    };

    return (
        <div className="workspace-shell">
            <aside className="catalog-panel">
                <div className="catalog-header">
                    <div className="catalog-toolbar">
                        <p className="catalog-kicker">{messages.app.labTitle}</p>
                        <div className="language-switch" aria-label={messages.app.languageLabel}>
                            {languageOptions.map((option) => (
                                <button
                                    key={option.value}
                                    className={`language-button ${language === option.value ? 'active' : ''}`}
                                    onClick={() => void handleLanguageChange(option.value)}
                                    type="button"
                                >
                                    {option.label}
                                </button>
                            ))}
                        </div>
                    </div>
                    <h1>{messages.app.workspaceTitle}</h1>
                    <p className="catalog-copy">{messages.app.workspaceSummary}</p>
                </div>

                <div className="catalog-list">
                    {demos.map((demo, index) => {
                        const active = demo.id === activeDemo;

                        return (
                            <button
                                key={demo.id}
                                className={`catalog-item ${active ? 'active' : ''}`}
                                onClick={() => setActiveDemo(demo.id)}
                                type="button"
                            >
                                <span className="catalog-index">0{index + 1}</span>
                                <div>
                                    <span className="catalog-badge">{demo.badge}</span>
                                    <strong>{demo.title}</strong>
                                    <p>{demo.summary}</p>
                                </div>
                            </button>
                        );
                    })}
                </div>
            </aside>

            <main className="demo-stage">
                <header className="stage-header">
                    <p className="stage-badge">{currentDemo.badge} Demo</p>
                    <h2>{currentDemo.title}</h2>
                    <p>{currentDemo.summary}</p>
                    <div className="stage-focus">
                        <span>{messages.app.focusLabel}</span>
                        <strong>{currentDemo.focus}</strong>
                    </div>
                </header>

                {activeDemo === 'counter'
                    ? <CounterDemoPanel messages={messages}/>
                    : <DownloadDemoPanel messages={messages}/>}
            </main>
        </div>
    );
}

function getStoredLanguage(): Language | null {
    if (typeof window === 'undefined') {
        return null;
    }

    const language = window.localStorage.getItem(languageStorageKey);
    return isLanguage(language) ? language : null;
}

function applyLanguage(language: Language, setLanguage: (language: Language) => void) {
    setLanguage(language);
    document.documentElement.lang = language;
    window.localStorage.setItem(languageStorageKey, language);
}

export default App;
