
import {EditorState} from "@codemirror/state"
import {EditorView, keymap, placeholder} from "@codemirror/view"
import {defaultKeymap} from "@codemirror/commands"

import './styles/index.scss';

let startState = EditorState.create({
    doc: "",
    extensions: [keymap.of(defaultKeymap), placeholder("# paste text here…")]
})

let view = new EditorView({
    state: startState,
    parent: document.getElementById("input-text-editor"),
})

view.focus();

