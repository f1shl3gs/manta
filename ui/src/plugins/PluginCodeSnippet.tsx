import React from 'react'
import CodeSnippet from '../shared/components/CodeSnippet'

interface Props {
  code: string
  language: string
}

const PluginCodeSnippet: React.FC<Props> = (props) => {
  const {code, language} = props

  return <CodeSnippet copyText={code} label={language} />
}

export default PluginCodeSnippet
