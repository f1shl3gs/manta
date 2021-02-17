import React from 'react'

interface Props {
  left?: JSX.Element[] | JSX.Element
  right?: JSX.Element[] | JSX.Element
}

const TabbedPageHeader: React.FC<Props> = ({left, right}) => {
  let leftHeader = <></>
  let rightHeader = <></>

  if (left) {
    leftHeader = <div className={'tabbed-page--header-left'}>{left}</div>
  }

  if (right) {
    rightHeader = <div className={'tabbed-page--header-right'}>{right}</div>
  }

  return (
    <div className="tabbed-page--header">
      {leftHeader}
      {rightHeader}
    </div>
  )
}

export default TabbedPageHeader
