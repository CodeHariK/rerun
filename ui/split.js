import Split from '/ui/cdn.split.js';

import { editor } from '/ui/monaco.js';

Split(['#treeview', '#codeeditor', '#outputWindow', '#consoleBox'], {
    sizes: [8, 36, 36, 20],
    minSize: 2,
    gutterSize: 3,
    cursor: 'col-resize',
    onDrag: function (sizes) {
        editor.layout()
    },
    onDragEnd: function (sizes) {
        localStorage.setItem('split-sizes', JSON.stringify(sizes))
    },
})

Split(['#pageUrl', '#contentFrame'], {
    direction: 'vertical',
    sizes: [5, 95],
    minSize: 2,
    gutterSize: 3,
    cursor: 'row-resize'
})

Split(['#statusPane', '#logsPane', "#terminal", "#commandBox"], {
    direction: 'vertical',
    sizes: [5, 70, 20, 5],
    minSize: 2,
    gutterSize: 3,
    cursor: 'row-resize'
})