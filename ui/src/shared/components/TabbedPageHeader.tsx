import React, {CSSProperties, FunctionComponent} from 'react'

interface Props {
  left: JSX.Element | JSX.Element[]
  right: JSX.Element | JSX.Element[]
  style?: CSSProperties
}

const TabbedPageHeader: FunctionComponent<Props> = ({left, right, style}) => {
  return (
    <div className={'tabbed-page--header'} style={style}>
      <div className={'tabbed-page--header-left'}>{left}</div>

      <div className={'tabbed-page--header-right'}>{right}</div>
    </div>
  )
}

export default TabbedPageHeader
