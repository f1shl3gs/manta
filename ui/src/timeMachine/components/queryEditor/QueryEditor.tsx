// Libraries
import React, {FunctionComponent} from 'react'

// Componetns
import {promLanguageDefinition} from 'monaco-promql'
import MonacoEditor from 'react-monaco-editor'
import PromqlToolbar from 'src/timeMachine/components/queryEditor/PromqlToolbar'
import ThemeName from 'src/shared/components/editor/PromQLEditorTheme'

// Types
import {EditorType} from 'src/types/monaco'

interface Props {
  query: string
  onChange: (text: string) => void
}

const QueryEditor: FunctionComponent<Props> = ({query, onChange}) => {
  const editorWillMount = monaco => {
    const languageId = promLanguageDefinition.id

    /*
      TODO: add hotkey for submit, e.g. Ctrl + Enter
    */

    monaco.languages.register(promLanguageDefinition)
    monaco.languages.onLanguage(languageId, () => {
      promLanguageDefinition.loader().then(mod => {
        monaco.languages.setMonarchTokensProvider(languageId, mod.language)
        monaco.languages.setLanguageConfiguration(
          languageId,
          mod.languageConfiguration
        )
        monaco.languages.registerCompletionItemProvider(
          languageId,
          mod.completionItemProvider
        )
      })
    })
  }

  const editorDidMount = (editor: EditorType) => {
    editor.focus()
  }

  return (
    <div className={'flux-editor'}>
      <div className={'flux-editor--left-panel'}>
        <MonacoEditor
          language="promql"
          theme={ThemeName}
          onChange={onChange}
          value={query}
          editorWillMount={editorWillMount}
          editorDidMount={editorDidMount}
        />
      </div>

      <div className={'flux-editor--right-panel'}>
        <PromqlToolbar />
      </div>
    </div>
  )
}

export default QueryEditor
