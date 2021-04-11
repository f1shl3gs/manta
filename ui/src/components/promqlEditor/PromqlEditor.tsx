// Libraries
import React, {useEffect, useRef} from 'react'

// Components
import {CompleteStrategy, newCompleteStrategy} from 'codemirror-promql/complete'
import {
  autocompletion,
  completionKeymap,
  CompletionContext,
  CompletionResult,
} from '@codemirror/autocomplete'
import {indentOnInput, syntaxTree} from '@codemirror/language'
import {history, historyKeymap} from '@codemirror/history'
import {
  EditorView,
  highlightSpecialChars,
  keymap,
  placeholder,
  ViewUpdate,
} from '@codemirror/view'
import {PromQLExtension} from 'codemirror-promql'
import {lineNumbers} from '@codemirror/gutter'
import {Compartment, EditorState, Prec} from '@codemirror/state'
import {bracketMatching} from '@codemirror/matchbrackets'
import {closeBrackets, closeBracketsKeymap} from '@codemirror/closebrackets'
import {defaultKeymap, insertNewlineAndIndent} from '@codemirror/commands'
import {commentKeymap} from '@codemirror/comment'
import {lintKeymap} from '@codemirror/lint'

import {promqlHighlighter, theme} from './theme'
import {useOrgID} from '../../shared/useOrg'

const promqlExtension = new PromQLExtension()
const dynamicConfigCompartment = new Compartment()

// Autocompletion strategy that wraps the main one and enriches
// it with past query items.
export class HistoryCompleteStrategy implements CompleteStrategy {
  private complete: CompleteStrategy
  private queryHistory: string[]

  constructor(complete: CompleteStrategy, queryHistory: string[]) {
    this.complete = complete
    this.queryHistory = queryHistory
  }

  promQL(
    context: CompletionContext
  ): Promise<CompletionResult | null> | CompletionResult | null {
    return Promise.resolve(this.complete.promQL(context)).then(res => {
      const {state, pos} = context
      const tree = syntaxTree(state).resolve(pos, -1)
      const start = res != null ? res.from : tree.from

      if (start !== 0) {
        return res
      }

      const historyItems: CompletionResult = {
        from: start,
        to: pos,
        options: this.queryHistory.map(q => ({
          label: q.length < 80 ? q : q.slice(0, 76).concat('...'),
          detail: 'past query',
          apply: q,
          info: q.length < 80 ? undefined : q,
        })),
        span: /^[a-zA-Z0-9_:]+$/,
      }

      if (res !== null) {
        historyItems.options = historyItems.options.concat(res.options)
      }
      return historyItems
    })
  }
}

interface Props {
  value: string
  onChange: (v: string) => void
}

const PromqlEditor: React.FC<Props> = props => {
  const {value, onChange} = props
  const containerRef = useRef<EditorView | null>(null)
  const viewRef = useRef<EditorView | null>(null)
  const orgID = useOrgID()
  const pathPrefix = `/api/v1/query/${orgID}`

  // settings
  const enableHighlighting = false
  const enableAutocomplete = true
  const enableLinter = false

  // props
  const queryHistory = ['']

  // (Re)initialize editor based on settings/setting changes
  useEffect(() => {
    // Build the dynamic part of the config
    promqlExtension.activateCompletion(enableAutocomplete)
    promqlExtension.activateLinter(enableLinter)
    promqlExtension.setComplete({
      completeStrategy: new HistoryCompleteStrategy(
        newCompleteStrategy({
          remote: {url: pathPrefix},
        }),
        queryHistory
      ),
    })
    const dynamicConfig = [
      enableHighlighting ? promqlHighlighter : [],
      promqlExtension.asExtension(),
    ]

    // Create or reconfigure the editor
    const view = viewRef.current
    if (view === null) {
      // If the editor does not exist yet, create one
      if (!containerRef.current) {
        throw new Error('expected CodeMirror container element to exist')
      }

      const startState = EditorState.create({
        doc: value,
        extensions: [
          theme,
          lineNumbers(),
          highlightSpecialChars(),
          history(),
          EditorState.allowMultipleSelections.of(true),
          indentOnInput(),
          bracketMatching(),
          closeBrackets(),
          autocompletion(),
          keymap.of([
            ...closeBracketsKeymap,
            ...defaultKeymap,
            // ...searchKeymap,
            ...historyKeymap,
            ...commentKeymap,
            ...completionKeymap,
            ...lintKeymap,
          ]),
          placeholder('Expression (press Shift+Enter for query)'),
          dynamicConfigCompartment.of(dynamicConfig),
          // This keymap is added without precedence so that closing the autocomplete dropdown
          // via Escape works without blurring the editor
          keymap.of([
            {
              key: 'Escape',
              run: (v: EditorView): boolean => {
                v.contentDOM.blur()
                return false
              },
            },
          ]),
          Prec.override(
            keymap.of([
              /*{
                key: 'Enter',
                run: insertNewlineAndIndent,
              },*/
              {
                key: 'Shift-Enter',
                run: (v: EditorView): boolean => {
                  console.log('do query')
                  return true
                },
              },
            ])
          ),
          EditorView.updateListener.of((upd: ViewUpdate): void => {
            // onChange(upd.state.doc.toString())
            const text = upd.state.doc.toString()
            // onChange(text)
          }),
        ],
      })

      const view = new EditorView({
        state: startState,
        // @ts-ignore
        parent: containerRef.current,
      })

      viewRef.current = view

      view.focus()
    } else {
      // The editor already exists, just reconfigure the dynamically configured parts
      view.dispatch(
        view.state.update({
          effects: dynamicConfigCompartment.reconfigure(dynamicConfig),
        })
      )
    }

    // value is only used in the initial render, so we don't want to
    // re-run this effect every time that value changes
    //
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [
    enableAutocomplete,
    enableHighlighting,
    enableLinter,
    queryHistory,
    onChange,
  ])

  return (
    <div
      // @ts-ignore
      ref={containerRef}
      className={'cm-expression-input'}
    />
  )
}

PromqlEditor.whyDidYouRender = true

export default PromqlEditor
