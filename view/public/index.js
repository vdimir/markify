
import {EditorState} from "@codemirror/state"
import {EditorView, keymap, placeholder} from "@codemirror/view"
import {defaultKeymap} from "@codemirror/commands"

import './styles/index.scss';


(function () {
    let editorRoot = document.getElementById("main-text-area")
    if (!editorRoot) {
        throw '"main-text-area" not found!'
    }
    let editorView = buildEditor(editorRoot);
    setupEventLisners(editorView)
})()

/// Replace html textarea with CodeMirror editor, setup button listeners
function buildEditor(editorRoot) {
    let originalTextField = document.getElementById("simple-text-area")

    let startState = EditorState.create({
        doc: originalTextField.textContent,
        extensions: [keymap.of(defaultKeymap), placeholder("# paste text here…")]
    })

    let view = new EditorView({
        state: startState,
        parent: editorRoot,
    })

    originalTextField.required = false
    originalTextField.hidden = true
    editorRoot.hidden = false

    let editorForm = document.getElementById("editor-form")
    editorForm.onsubmit = onDataSubmit(view)

    view.focus();
    
    return view;
}

function setupEventLisners(editorView) {
    document.getElementsByClassName('alert-message')[0].addEventListener('click', () => showAlert(false) )
    document.getElementsByClassName('btn-preview-text')[0].addEventListener('click', () => previewToggle(editorView) )
}

function showAlert(msg) {
    const boxClassName = "editor-top-box";
    const contentClassName = "alert-message";
    if (msg === false) {
        document.getElementsByClassName(boxClassName)[0].classList.remove('visible')
        return;
    }

    document.getElementsByClassName(contentClassName)[0].innerHTML = msg;
    document.getElementsByClassName(boxClassName)[0].classList.add('visible')
}

function formRequest(data) {
    // TODO: How to stream data without copying?
    return {
        body: "data=" + data,
        headers: { "Content-Type": "application/x-www-form-urlencoded" },
        method: "POST"
    }
}

function validateInput(text) {
    if (!/\S/.test(text)) {
        showAlert("Insert some text!")
        return false
    }
    return true
}

function displayPreview(resp) {
    let text = resp['body']
    document.getElementById("root-textarea").style.display = 'none'
    
    let contentRoot = document.getElementById('preview-content');
    contentRoot.innerHTML = text
    contentRoot.style.display = 'block';
    
    document.querySelector('.btn-preview-text > .fa').classList.add('fa-eye-slash')
}

function hidePreview(editorView) {
    let contentRoot = document.getElementById('preview-content');
    contentRoot.style.display = 'none';
    
    document.getElementById("root-textarea").style.display = 'block'
    editorView.focus()
    
    document.querySelector('.btn-preview-text > .fa').classList.remove('fa-eye-slash')
}

function isInPreview() {
    let contentRoot = document.getElementById('preview-content');
    return contentRoot.style.display != 'none';
}

function apiCall(path, body, callback) {
    fetch(path, {
        headers: { "Content-Type": "application/json" },
        method: "POST",
        body: JSON.stringify(body),
    }).then(resp => {
        if (resp.ok) {
            return resp.json()
        }
        throw resp.statusText
    })
    .then(callback)
    .catch(err => {
        showAlert(err)
    });
}

function previewText(editorView) {
    let data = editorView.state.doc.toString();
    
    if (!validateInput(data)) {
        return
    }
    
    apiCall("/api/preview", {'text': data}, displayPreview)
}

function previewToggle(editorView) {
    if (isInPreview()) {
        hidePreview(editorView)
    } else {
        previewText(editorView)
    }
}

function onDataSubmit(view) {
    return e => {
        console.log(e)
        e.preventDefault()
        let data = view.state.doc.toString();
        if (!validateInput(data)) {
            return false
        }
        apiCall("/api/create", {'text': data}, resp => { document.location.href = resp['path'] })
        return false
    }
}
