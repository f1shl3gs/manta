// Libraries
import React, {FunctionComponent} from 'react'
import {connect, ConnectedProps} from 'react-redux'

// Componetns
import {promLanguageDefinition} from 'monaco-promql'
import MonacoEditor from 'react-monaco-editor'
import PromqlToolbar from 'src/timeMachine/components/PromqlToolbar'
import ThemeName from 'src/shared/components/editor/PromQLEditorTheme'

// Types
import {EditorType} from 'src/types/monaco'
import {AppState} from 'src/types/stores'

// Actions
import {setActiveQueryText} from 'src/timeMachine/actions'

const mstp = (state: AppState) => {
  const {viewProperties, activeQueryIndex} = state.timeMachine

  return {
    query: viewProperties.queries[activeQueryIndex].text,
  }
}

const mdtp = {
  onChange: setActiveQueryText,
}

const connector = connect(mstp, mdtp)

type Props = ConnectedProps<typeof connector>

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

export default connector(QueryEditor)
