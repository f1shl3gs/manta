// Libraries
import React, {CSSProperties} from 'react'
import {DapperScrollbars} from '@influxdata/clockface'
import classnames from 'classnames'

interface Props {
  scrollable?: boolean
  addPadding: boolean
  testID?: string
  autoHideScrollbars: boolean
  style?: CSSProperties
  className?: string
}

const BuilderCardBody: React.FC<Props> = props => {
  const {
    style,
    className,
    scrollable = true,
    addPadding = true,
    testID = 'builder-card--body',
    autoHideScrollbars = false,
  } = props

  const children = addPadding ? (
    <div className={'builder-card--contents'}>{props.children}</div>
  ) : (
    props.children
  )

  if (scrollable) {
    const scrollbarStyles = {maxWidth: '100%', maxHeight: '100%', ...style}

    return (
      <DapperScrollbars
        className={'builder-card--body'}
        style={scrollbarStyles}
        testID={testID}
        autoHide={autoHideScrollbars}
      >
        {children}
      </DapperScrollbars>
    )
  }

  const classname = classnames('builder-card--body', {
    [`${className}`]: className,
  })

  return (
    <div className={classname} data-testid={testID} style={style}>
      {children}
    </div>
  )
}

export default BuilderCardBody
