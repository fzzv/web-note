import {useState} from 'react';
import './App.css';
import {CounterDemoPanel} from './demos/CounterDemoPanel';
import {DownloadDemoPanel} from './demos/DownloadDemoPanel';

type DemoId = 'counter' | 'download';

type DemoDefinition = {
    id: DemoId;
    badge: string;
    title: string;
    summary: string;
    focus: string;
};

const demos: DemoDefinition[] = [
    {
        id: 'counter',
        badge: 'Lifecycle',
        title: 'Persistent Counter',
        summary: 'Keep the original counter demo, but isolate it as a reusable module with startup restore and shutdown persistence.',
        focus: 'Wails lifecycle hooks + local state persistence',
    },
    {
        id: 'download',
        badge: 'Events',
        title: 'Download Progress Notifications',
        summary: 'Stream a file from Go, emit progress events through Wails, and render the download state live in React.',
        focus: 'runtime.EventsEmit / EventsOn + async backend work',
    },
];

function App() {
    const [activeDemo, setActiveDemo] = useState<DemoId>('counter');

    const currentDemo = demos.find((demo) => demo.id === activeDemo) ?? demos[0];

    return (
        <div className="workspace-shell">
            <aside className="catalog-panel">
                <div className="catalog-header">
                    <p className="catalog-kicker">Wails Demo Lab</p>
                    <h1>Multiple examples in one project.</h1>
                    <p className="catalog-copy">
                        Use the project as a small demo workspace: each sample keeps its own backend module and frontend panel.
                    </p>
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
                        <span>Learning focus</span>
                        <strong>{currentDemo.focus}</strong>
                    </div>
                </header>

                {activeDemo === 'counter' ? <CounterDemoPanel/> : <DownloadDemoPanel/>}
            </main>
        </div>
    );
}

export default App;
