// Libraries
import React from 'react'
import ReactMarkdown from 'react-markdown'

interface Props {
  className?: string
  text: string
}

const MarkdownRenderer: React.FC<Props> = ({className, text}) => {
  return <ReactMarkdown source={text} className={className} />
}

export default MarkdownRenderer
