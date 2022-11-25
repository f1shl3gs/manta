import {monaco} from 'react-monaco-editor'
import {MonacoType} from 'src/types/monaco'

const ThemeName = 'yaml'

function addTheme(monaco: MonacoType) {
  monaco.editor.defineTheme(ThemeName, {
    base: 'vs-dark',
    inherit: false,

    // keys are found in
    // https://github.com/microsoft/monaco-editor/blob/136ce723f73b8bd284565c0b7d6d851b52161015/src/basic-languages/yaml/yaml.ts
    rules: [
      {
        token: 'comment',
        foreground: '#7CE490',
      },
      {
        token: 'string',
        foreground: '#7CE490',
      },

      // numbers
      {
        token: 'number',
        foreground: '#7CE490',
      },

      // others
      {
        token: '',
        foreground: '#f8f8f8',
        background: '#202028',
      },
    ],
    colors: {
      'editor.foreground': '#F8F8F8',
      'editor.background': '#202028',
      'editorGutter.background': '#25252e',
      'editor.selectionBackground': '#353640',
      'editorLineNumber.foreground': '#666978',
      'editor.lineHighlightBackground': '#353640',
      'editorCursor.foreground': '#ffffff',
      'editorActiveLineNumber.foreground': '#bec2cc',
    },
  })
}

addTheme(monaco)

export default ThemeName
