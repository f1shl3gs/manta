// Libraries
import React, {FunctionComponent} from 'react'

// Components
import MonacoEditor from 'react-monaco-editor'
import {EditorType} from 'src/types/monaco'
import ThemeName from 'src/shared/components/YamlEditorTheme'

interface Props {
  content: string
  testID?: string
  readOnly?: boolean
  onChange?: (content: string) => void
}

const YamlMonacoEditor: FunctionComponent<Props> = props => {
  const {content, onChange, readOnly, testID = 'yaml-editor'} = props

  const editorDidMount = (editor: EditorType) => {
    editor.focus()
  }

  return (
    <div data-testid={testID} className={'editor--embeded'}>
      <MonacoEditor
        language={'yaml'}
        theme={ThemeName}
        value={content}
        onChange={onChange}
        editorDidMount={editorDidMount}
        options={{
          fontFamily: '"IBMPlexMono", monospace',
          cursorWidth: 2,
          tabSize: 2,
          lineNumbersMinChars: 4,
          lineDecorationsWidth: 0,
          minimap: {
            renderCharacters: false,
          },
          overviewRulerBorder: false,
          automaticLayout: true,
          readOnly: readOnly || false,
        }}
      />
    </div>
  )
}

export default YamlMonacoEditor
