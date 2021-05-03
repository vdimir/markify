
import {EditorState} from '@codemirror/state'
import {EditorView, keymap, placeholder} from '@codemirror/view'
import {defaultKeymap} from '@codemirror/commands'

import './styles/index.scss';


(function () {
    let editorRoot = document.getElementById('main-text-area')
    if (!editorRoot) {
        throw 'editor root not found'
    }
    let editorView = buildEditor(editorRoot);
    setupEventLisners(editorView)
})()

/// Replace html textarea with CodeMirror editor, setup button listeners
function buildEditor(editorRoot) {
    let originalTextField = document.getElementById('simple-text-area')

    let startState = EditorState.create({
        doc: originalTextField.textContent,
        extensions: [keymap.of(defaultKeymap), placeholder('# paste text here…')]
    })

    let view = new EditorView({
        state: startState,
        parent: editorRoot,
    })

    originalTextField.required = false
    originalTextField.hidden = true
    editorRoot.hidden = false


    view.focus();
    
    return view;
}

function setupEventLisners(editorView) {
    let editorForm = document.getElementById('editor-form')
    editorForm.onsubmit = e => { onDataSubmit(e, editorView).catch(showAlert) }

    document.getElementsByClassName('alert-close')[0].addEventListener('click', () => showAlert(false) )
    document.getElementsByClassName('btn-preview-text')[0].addEventListener('click', () => previewToggle(editorView).catch(showAlert) )
    document.getElementById('preview-content').addEventListener('dblclick', () => previewToggle(editorView).catch(showAlert) )
}

function showAlert(msg) {
    const boxClassName = 'editor-top-box';
    const contentClassName = 'alert-message';
    if (msg === false) {
        document.getElementsByClassName(boxClassName)[0].classList.remove('visible')
        return;
    }

    document.getElementsByClassName(contentClassName)[0].innerHTML = msg;
    document.getElementsByClassName(boxClassName)[0].classList.add('visible')
}

function validateInput(text) {
    if (!/\S/.test(text)) {
        throw 'Insert some text!'
    }
}

function displayPreview(resp) {
    let text = resp['body']
    document.getElementById('root-textarea').style.display = 'none'
    
    let contentRoot = document.getElementById('preview-content');
    contentRoot.innerHTML = text
    contentRoot.style.display = 'block';
    
    document.querySelector('.btn-preview-text > .fa').classList.add('fa-eye-slash')
}

async function hidePreview(editorView) {
    let contentRoot = document.getElementById('preview-content');
    contentRoot.style.display = 'none';
    
    document.getElementById('root-textarea').style.display = 'block'
    editorView.focus()
    
    document.querySelector('.btn-preview-text > .fa').classList.remove('fa-eye-slash')
}

function isInPreview() {
    let contentRoot = document.getElementById('preview-content');
    return contentRoot.style.display != 'none';
}

async function apiCall(path, body) {
    return fetch(path, {
        headers: { 'Content-Type': 'application/json' },
        method: 'POST',
        body: JSON.stringify(body),
    }).then(resp => {
        if (resp.ok) {
            return resp.json()
        }
        throw resp.statusText
    })
}

async function previewText(editorView) {
    let data = editorView.state.doc.toString();
    validateInput(data)
    
    displayPreview(await apiCall('/api/preview', {'text': data}))
}

async function previewToggle(editorView) {
    if (isInPreview()) {
        await hidePreview(editorView)
    } else {
        await previewText(editorView)
    }
}

async function onDataSubmit(event, view) {
    event.preventDefault()
    let data = view.state.doc.toString();
    validateInput(data)

    let resp = await apiCall('/api/create', {'text': data})
    document.location.href = resp['path']
    return false
}
