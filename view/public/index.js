
import {EditorState} from "@codemirror/state"
import {EditorView, keymap, placeholder} from "@codemirror/view"
import {defaultKeymap} from "@codemirror/commands"

import './styles/index.scss';


(function () {
    let editorRoot = document.getElementById("main-text-area")
    if (!editorRoot) {
        throw '"main-text-area" not found!'
    }
    buildEditor(editorRoot)
    setupEventLisners()
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
}

function setupEventLisners() {
    document.getElementsByClassName('alert-message')[0].addEventListener('click', () => showAlert(false) )
}

export function showAlert(msg) {
    const boxClassName = "editor-top-box";
    const contentClassName = "alert-message";
    if (msg === false) {
        document.getElementsByClassName(boxClassName)[0].classList.remove('visible')
        return;
    }

    document.getElementsByClassName(contentClassName)[0].innerHTML = msg;
    document.getElementsByClassName(boxClassName)[0].classList.add('visible')
}

function onDataSubmit(view) {
    return e => {
        e.preventDefault()
        let data = view.state.doc.toString();
        if (!/\S/.test(data)) {
            showAlert("Insert some text!")
            return false
        }
        fetch("/create", {
            body: "data=" + data,
            headers: { "Content-Type": "application/x-www-form-urlencoded" },
            method: "POST"
        }).then(resp => {
            if (resp.ok) {
                document.location.href = resp.url
                return
            }
            throw resp.statusText
        }).catch(_error => {
            showAlert("Something went wrong")
        });
        return false
    }
}
