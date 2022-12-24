// Libraries
import React, {CSSProperties, FunctionComponent, ReactNode} from 'react'
import classnames from 'classnames'

// Components
import {ComponentSize, DapperScrollbars} from '@influxdata/clockface'

interface Props {
  scrollable?: boolean
  autoHideScrollbars: boolean
  addPadding: boolean
  testID?: string
  children: JSX.Element | JSX.Element[]
  className?: string
  style?: CSSProperties
}

const BuilderCardBody: FunctionComponent<Props> = ({
  scrollable = true,
  autoHideScrollbars = false,
  addPadding,
  testID = 'builder-card--body',
  className,
  style,
  children,
}) => {
  const content = (): JSX.Element | ReactNode => {
    if (addPadding) {
      return <div className={'builder-card--contents'}>{children}</div>
    }

    return children
  }

  if (scrollable) {
    const scrollbarStyles = {
      maxWidth: '100%',
      maxHeight: '100%',
      ...style,
    }

    return (
      <DapperScrollbars
        className={'builder-card--body'}
        style={scrollbarStyles}
        testID={testID}
        size={ComponentSize.ExtraSmall}
        autoHide={autoHideScrollbars}
      >
        {content()}
      </DapperScrollbars>
    )
  }

  const classname = classnames('builder-card--body', {
    [`${className}`]: className,
  })

  return (
    <div className={classname} style={style} data-testid={testID}>
      {content()}
    </div>
  )
}

export default BuilderCardBody
