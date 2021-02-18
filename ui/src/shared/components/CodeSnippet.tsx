// Libraries
import React from 'react'
import {DapperScrollbars} from '@influxdata/clockface'
import CopyButton from './CopyButton'

interface Props {
  copyText: string
  label: string
  onClick?: () => void
}

const CodeSnippet: React.FC<Props> = props => {
  const {copyText, label, onClick} = props

  return (
    <div className={'code-snippet'}>
      <DapperScrollbars
        autoHide={false}
        autoSizeHeight={true}
        className={'code-snippet--scroll'}
      >
        <div className={'code-snippet--text'}>
          <pre>
            <code>{copyText}</code>
          </pre>
        </div>
      </DapperScrollbars>

      <div className={'code-snippet--footer'}>
        <CopyButton
          textToCopy={copyText}
          //@ts-ignore
          onCopyText={(text, status) => console.log('onCopyText', text, status)}
          contentName={'Script'}
          onClick={onClick}
          buttonText={'Copy to Clipboard'}
        />

        <label className={'code-snippet--label'}>{label}</label>
      </div>
    </div>
  )
}

export default CodeSnippet
