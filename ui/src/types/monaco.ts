import * as allMonaco from 'monaco-editor/esm/vs/editor/editor.api'

export type MonacoType = typeof allMonaco
export type EditorType = allMonaco.editor.IStandaloneCodeEditor
export type CursorEvent = allMonaco.editor.ICursorPositionChangedEvent
export type KeyboardEvent = allMonaco.IKeyboardEvent
