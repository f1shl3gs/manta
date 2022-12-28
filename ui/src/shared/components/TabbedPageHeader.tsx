import {FunctionComponent} from 'react'

interface Props {
  left: JSX.Element | JSX.Element[]
  right: JSX.Element | JSX.Element[]
}

const TabbedPageHeader: FunctionComponent<Props> = ({left, right}) => {
  return (
    <div className={'tabbed-page--header'}>
      <div className={'tabbed-page--header-left'}>{left}</div>

      <div className={'tabbed-page--header-right'}>{right}</div>
    </div>
  )
}

export default TabbedPageHeader
