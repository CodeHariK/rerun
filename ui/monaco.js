import * as monaco from '/ui/cdn.monaco.js';

let lastFileSaveTime;

var languageMappings = {
    '.js': 'javascript',
    '.ts': 'typescript',
    '.py': 'python',
    '.java': 'java',
    '.cpp': 'cpp',
    '.php': 'php',
    '.go': 'go',
    '.html': 'html',
    '.css': 'css',
    '.json': 'json',
    '.yaml': 'yaml',
    '.xml': 'xml',
    '.md': 'markdown',
};

export const editor = monaco.editor.create(document.getElementById('codeeditor'), {
    value: "",
    language: 'go',
    theme: 'vs',
    minimap: {
        enabled: false
    },
});

export function monacoToggleTheme(dark) {
    if (dark) {
        monaco.editor.setTheme("vs-dark")
    } else {
        monaco.editor.setTheme("vs")
    }
}

setTimeout(function () {
    editor.layout()
}, 100)
window.addEventListener('resize', function handleResize() {
    setTimeout(function () {
        editor.layout()
    }, 100)
});

// Event listener for keydown events
document.addEventListener('keydown', function handleSaveShortcut(event) {
    if ((event.ctrlKey || event.metaKey) && event.key === 's') {
        event.preventDefault();
        saveContent();
    }
});

export async function fetchFile(filePath) {
    const url = `http://localhost:9753/file?filepath=${encodeURIComponent(filePath)}`;

    try {
        const response = await fetch(url);

        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }

        const data = await response.text();

        if (editor.getValue() == data) {
            return
        }

        editor.setValue(data);
        var fileExtension = filePath.substr(filePath.lastIndexOf('.')).toLowerCase();
        if (languageMappings[fileExtension]) {
            monaco.editor.setModelLanguage(editor.getModel(), languageMappings[fileExtension]);
        } else {
            console.log('Unsupported file extension or language');
        }

    } catch (error) {
        console.error('Error fetching file:', error);
    }
}

export async function saveContent() {
    if (lastFileSaveTime + 1000 > Date.now()) {
        return
    }
    lastFileSaveTime = Date.now()

    const content = editor.getValue();

    if (!window.currentFile) {
        return;
    }

    const response = await fetch('http://localhost:9753/save', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ currentFile: window.currentFile, content })
    });

    if (response.ok) {
        console.log('Content saved successfully');
    } else {
        console.log('Error saving content');
    }
}
window.saveContent = saveContent;

window.currentFile = localStorage.getItem("currentFile")
fetchFile(window.currentFile)