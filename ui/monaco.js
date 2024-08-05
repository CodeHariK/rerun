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

monaco.editor.defineTheme('atom-dark', {
    base: 'vs-dark', // Use the vs-dark base theme for Monaco
    inherit: true, // Inherit from the base theme
    rules: [
        // Comments and Quotes
        { token: 'comment', foreground: '5c6370', fontStyle: 'italic' },
        { token: 'quote', foreground: '5c6370', fontStyle: 'italic' },

        // Keywords and Doctags
        { token: 'keyword', foreground: 'c678dd' },
        { token: 'keyword.control', foreground: 'c678dd' },
        { token: 'keyword.operator', foreground: 'c678dd' },
        { token: 'keyword.other', foreground: 'c678dd' },

        // Sections, Names, Tags
        { token: 'variable.language', foreground: 'e06c75' },
        { token: 'variable.other.local', foreground: 'e06c75' },
        { token: 'variable.other', foreground: 'e06c75' },
        { token: 'type', foreground: 'e06c75' },
        { token: 'function', foreground: 'e06c75' },

        // Literals
        { token: 'number', foreground: '56b6c2' },
        { token: 'string', foreground: '98c379' },
        { token: 'regexp', foreground: '98c379' },
        { token: 'literal', foreground: '56b6c2' },

        // Built-ins and Classes
        { token: 'type', foreground: 'e6c07b' },
        { token: 'class', foreground: 'e6c07b' },
        { token: 'variable', foreground: 'd19a66' },

        // Attributes, Variables, Template Variables
        { token: 'attribute', foreground: 'd19a66' },
        { token: 'variable', foreground: 'd19a66' },
        { token: 'variable.parameter', foreground: 'd19a66' },

        // Symbols, Links, Titles
        { token: 'symbol', foreground: '61aeee' },
        { token: 'link', foreground: '61aeee', fontStyle: 'underline' },
        { token: 'title', foreground: '61aeee' },

        // Emphasis and Strong
        { token: 'emphasis', fontStyle: 'italic' },
        { token: 'strong', fontStyle: 'bold' },
    ],
    colors: {
        'editor.background': '#282c34',
        'editor.foreground': '#abb2bf',
        'editorCursor.foreground': '#528b8b',
        'editor.lineHighlightBackground': '#2c313c',
        'editor.selectionBackground': '#3e4451',
        'editor.selectionHighlightBackground': '#3e4451',
        'editor.inactiveSelectionBackground': '#3e4451',
        'editorIndentGuide.background': '#3b4048',
        'editorIndentGuide.activeBackground': '#4b5363',
        'editorWhitespace.foreground': '#3b4048',
        'editorRuler.foreground': '#2c313c',
        'editorGroupHeader.tabsBackground': '#1c1e22',
        'editorGroupHeader.border': '#1c1e22',
        'tab.activeBackground': '#1c1e22',
        'tab.inactiveBackground': '#282c34',
        'tab.activeBorder': '#61aeee',
        'tab.unfocusedActiveBorder': '#61aeee',
    }
});

monaco.editor.defineTheme('atom-light', {
    base: 'vs', // Use the base 'vs' theme for Monaco (light theme)
    inherit: true, // Inherit from the base theme
    rules: [
        // Comments and Quotes
        { token: 'comment', foreground: 'a0a1a7', fontStyle: 'italic' },
        { token: 'quote', foreground: 'a0a1a7', fontStyle: 'italic' },

        // Keywords and Doctags
        { token: 'keyword', foreground: 'a626a4' },
        { token: 'keyword.control', foreground: 'a626a4' },
        { token: 'keyword.operator', foreground: 'a626a4' },
        { token: 'keyword.other', foreground: 'a626a4' },

        // Sections, Names, Tags
        { token: 'variable.language', foreground: 'e45649' },
        { token: 'variable.other.local', foreground: 'e45649' },
        { token: 'variable.other', foreground: 'e45649' },
        { token: 'type', foreground: 'e45649' },
        { token: 'function', foreground: 'e45649' },

        // Literals
        { token: 'number', foreground: '0184bb' },
        { token: 'string', foreground: '50a14f' },
        { token: 'regexp', foreground: '50a14f' },
        { token: 'literal', foreground: '0184bb' },

        // Built-ins and Classes
        { token: 'type', foreground: 'c18401' },
        { token: 'class', foreground: 'c18401' },
        { token: 'variable', foreground: '986801' },

        // Attributes, Variables, Template Variables
        { token: 'attribute', foreground: '986801' },
        { token: 'variable', foreground: '986801' },
        { token: 'variable.parameter', foreground: '986801' },

        // Symbols, Links, Titles
        { token: 'symbol', foreground: '4078f2' },
        { token: 'link', foreground: '4078f2', fontStyle: 'underline' },
        { token: 'title', foreground: '4078f2' },

        // Emphasis and Strong
        { token: 'emphasis', fontStyle: 'italic' },
        { token: 'strong', fontStyle: 'bold' },
    ],
    colors: {
        'editor.background': '#fafafa',
        'editor.foreground': '#383a42',
        'editorCursor.foreground': '#333333',
        'editor.lineHighlightBackground': '#f0f0f0',
        'editor.selectionBackground': '#e0e0e0',
        'editor.selectionHighlightBackground': '#e0e0e0',
        'editor.inactiveSelectionBackground': '#e0e0e0',
        'editorIndentGuide.background': '#dcdcdc',
        'editorIndentGuide.activeBackground': '#c0c0c0',
        'editorWhitespace.foreground': '#dcdcdc',
        'editorRuler.foreground': '#e0e0e0',
        'editorGroupHeader.tabsBackground': '#f5f5f5',
        'editorGroupHeader.border': '#f5f5f5',
        'tab.activeBackground': '#f5f5f5',
        'tab.inactiveBackground': '#fafafa',
        'tab.activeBorder': '#4078f2',
        'tab.unfocusedActiveBorder': '#4078f2',
    }
});

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
        // monaco.editor.setTheme("vs-dark")
        monaco.editor.setTheme("atom-dark")
    } else {
        // monaco.editor.setTheme("vs")
        monaco.editor.setTheme("atom-light")
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