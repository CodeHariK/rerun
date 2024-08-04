import { fetchFile } from '/ui/monaco.js';

export function generateFileTreeHTML(node, path = '') {
    const ul = document.createElement('ul');
    const currentPath = path ? `${path}/${node.name}` : node.name;

    if (node.is_dir) {
        const li = document.createElement('li');
        li.innerHTML = `<span class="folder-name" data-open="true">&#11153; ${node.name}</span>`;
        ul.appendChild(li);

        const childrenUl = document.createElement('ul');
        node.children.forEach(child => {
            childrenUl.appendChild(generateFileTreeHTML(child, currentPath));
        });
        li.appendChild(childrenUl);

        li.querySelector('span.folder-name').addEventListener('click', function () {
            const isOpen = this.getAttribute('data-open') === 'true';
            this.innerHTML = isOpen ? `&#9654; ${node.name}` : `&#11153; ${node.name}`;
            this.setAttribute('data-open', !isOpen);
            li.querySelector('ul').style.display = isOpen ? 'none' : 'block';
        });
    } else {
        const li = document.createElement('li');
        li.innerHTML = `<span>${node.name}</span>`;
        li.querySelector('span').addEventListener('click', () => {
            fetchFile(currentPath)
            window.currentFile = currentPath
            localStorage.setItem("currentFile", currentPath)
            console.log(currentPath);
        });
        ul.appendChild(li);
    }

    return ul;
}

export function fileTree(elementId) {
    NodeList.prototype.has = function (selector) {
        return Array.from(this).filter(e => e.querySelector(selector));
    };

    var element = document.getElementById(elementId);
    element.classList.add('file-list');
    var liElementsInideUl = element.querySelectorAll('li');
    liElementsInideUl.has('ul').forEach(li => {
        li.classList.add('folder-root', 'closed', 'open');
        var spanFolderElementsInsideLi = li.querySelectorAll('span.folder-name');
        spanFolderElementsInsideLi.forEach(span => {
            if (span.parentNode.nodeName === 'LI') {
                span.onclick = function (e) {
                    span.parentNode.classList.toggle('open');
                };
            }
        });
    });
}