export type Language = 'zh-CN' | 'en-US';

export type Messages = {
    app: {
        labTitle: string;
        workspaceTitle: string;
        workspaceSummary: string;
        focusLabel: string;
        languageLabel: string;
    };
    demos: {
        counter: {
            badge: string;
            title: string;
            summary: string;
            focus: string;
        };
        download: {
            badge: string;
            title: string;
            summary: string;
            focus: string;
        };
    };
    counterPanel: {
        loading: string;
        kicker: string;
        title: string;
        description: string;
        countLabel: string;
        decrement: string;
        increment: string;
        startupLabel: string;
        startupDescription: string;
        beforeCloseLabel: string;
        beforeCloseDescription: string;
        shutdownLabel: string;
        shutdownDescription: string;
    };
    downloadPanel: {
        loading: string;
        kicker: string;
        title: string;
        description: string;
        urlLabel: string;
        urlPlaceholder: string;
        note: string;
        start: string;
        downloading: string;
        reset: string;
        statusLabels: Record<string, string>;
        fileLabel: string;
        waitingFile: string;
        transferredLabel: string;
        savedToLabel: string;
        savedToPlaceholder: string;
        hints: string[];
    };
};

const translations: Record<Language, Messages> = {
    'zh-CN': {
        app: {
            labTitle: 'Wails 示例集',
            workspaceTitle: '一个项目，容纳多个 demo。',
            workspaceSummary: '把项目当成一个小型示例工作台：每个示例都拥有独立的后端模块和前端面板。',
            focusLabel: '学习重点',
            languageLabel: '语言',
        },
        demos: {
            counter: {
                badge: '生命周期',
                title: '持久化计数器',
                summary: '保留原有计数器示例，并把它拆成可复用模块，继续演示启动恢复和退出持久化。',
                focus: 'Wails 生命周期钩子 + 本地状态持久化',
            },
            download: {
                badge: '事件系统',
                title: '下载进度通知',
                summary: '由 Go 端发起文件下载，通过 Wails 事件持续推送进度，再由 React 实时渲染。',
                focus: 'runtime.EventsEmit / EventsOn + 异步后端任务',
            },
        },
        counterPanel: {
            loading: '正在加载计数器状态...',
            kicker: '独立模块',
            title: '计数器示例继续保留。',
            description: '计数器逻辑已经迁移到独立后端文件，但仍会在应用启动时恢复状态，并在退出时写回本地 JSON。',
            countLabel: '当前计数',
            decrement: '减一',
            increment: '加一',
            startupLabel: '启动阶段',
            startupDescription: '在 React 界面渲染前，从本地 JSON 文件中读取上次保存的计数。',
            beforeCloseLabel: '关闭前',
            beforeCloseDescription: '弹出确认框，方便观察 Wails 生命周期钩子的执行时机。',
            shutdownLabel: '退出阶段',
            shutdownDescription: '自动保存当前计数，下次重新打开应用时会继续沿用。',
        },
        downloadPanel: {
            loading: '正在加载下载示例...',
            kicker: 'runtime.EventsEmit',
            title: '观察后端进度事件如何推送到前端。',
            description: '点击开始下载后，方法会立即返回。Go 后端继续在 goroutine 中拉取响应、写入本地文件，并通过 <code>EventsOn</code> 持续推送状态。',
            urlLabel: '下载地址',
            urlPlaceholder: 'https://example.com/file.zip',
            note: '文件会保存到本机的 <code>Downloads/wails-demo</code> 目录。可填写任意可访问的 HTTP 或 HTTPS 文件地址。',
            start: '开始下载',
            downloading: '下载中...',
            reset: '重置状态',
            statusLabels: {
                idle: '空闲',
                starting: '准备中',
                downloading: '下载中',
                completed: '已完成',
                error: '失败',
            },
            fileLabel: '文件名',
            waitingFile: '等待新的下载任务',
            transferredLabel: '已传输',
            savedToLabel: '保存位置',
            savedToPlaceholder: 'Downloads/wails-demo',
            hints: [
                'Go 在 goroutine 中执行下载，因此不会阻塞前端界面线程。',
                '每一次事件都会携带完整的 DownloadState，前端可以直接渲染。',
                '下载成功或失败时，后端都会再发出一次终态事件。',
            ],
        },
    },
    'en-US': {
        app: {
            labTitle: 'Wails Demo Lab',
            workspaceTitle: 'Multiple demos in one project.',
            workspaceSummary: 'Treat the project as a small demo workspace: each sample has its own backend module and frontend panel.',
            focusLabel: 'Learning focus',
            languageLabel: 'Language',
        },
        demos: {
            counter: {
                badge: 'Lifecycle',
                title: 'Persistent Counter',
                summary: 'Keep the original counter sample and extract it into a reusable module while preserving startup restore and shutdown persistence.',
                focus: 'Wails lifecycle hooks + local state persistence',
            },
            download: {
                badge: 'Events',
                title: 'Download Progress Notifications',
                summary: 'Let Go download a file, emit progress through Wails events, and render the live state in React.',
                focus: 'runtime.EventsEmit / EventsOn + async backend work',
            },
        },
        counterPanel: {
            loading: 'Loading counter state...',
            kicker: 'Isolated module',
            title: 'The counter demo is still here.',
            description: 'The counter logic now lives in its own backend file, while still restoring state during startup and persisting JSON during shutdown.',
            countLabel: 'Current count',
            decrement: 'Decrement',
            increment: 'Increment',
            startupLabel: 'Startup',
            startupDescription: 'Reads the previously saved count from a local JSON file before the React UI renders.',
            beforeCloseLabel: 'Before close',
            beforeCloseDescription: 'Shows a confirmation dialog so the Wails lifecycle timing stays easy to observe.',
            shutdownLabel: 'Shutdown',
            shutdownDescription: 'Automatically saves the current count so reopening the app resumes the previous value.',
        },
        downloadPanel: {
            loading: 'Loading download demo...',
            kicker: 'runtime.EventsEmit',
            title: 'Watch backend progress events reach the frontend.',
            description: 'Starting a download returns immediately. The Go backend continues inside a goroutine, writes the file locally, and pushes state updates that React receives through <code>EventsOn</code>.',
            urlLabel: 'Download URL',
            urlPlaceholder: 'https://example.com/file.zip',
            note: 'Files are saved into your local <code>Downloads/wails-demo</code> folder. Use any reachable HTTP or HTTPS file URL.',
            start: 'Start download',
            downloading: 'Downloading...',
            reset: 'Reset status',
            statusLabels: {
                idle: 'Idle',
                starting: 'Starting',
                downloading: 'Downloading',
                completed: 'Completed',
                error: 'Error',
            },
            fileLabel: 'File',
            waitingFile: 'Waiting for a new download',
            transferredLabel: 'Transferred',
            savedToLabel: 'Saved to',
            savedToPlaceholder: 'Downloads/wails-demo',
            hints: [
                'Go starts the request in a goroutine, so the UI thread never blocks.',
                'Each event carries a complete DownloadState payload that React can render directly.',
                'When the request completes or fails, the backend emits one final terminal state.',
            ],
        },
    },
};

export const languageOptions: Array<{value: Language; label: string}> = [
    {value: 'zh-CN', label: '中文'},
    {value: 'en-US', label: 'English'},
];

export function getMessages(language: Language): Messages {
    return translations[language] ?? translations['zh-CN'];
}

export function normalizeLanguage(language: string | null | undefined): Language {
    const normalized = language?.trim().toLowerCase() ?? '';

    if (normalized.startsWith('en')) {
        return 'en-US';
    }

    return 'zh-CN';
}

export function isLanguage(language: string | null | undefined): language is Language {
    return language === 'zh-CN' || language === 'en-US';
}
